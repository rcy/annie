package image

import (
	"context"
	"fmt"
	"goirc/db/model"
	db "goirc/model"
	"io"
	"log"
	"net/http"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func mustGetenv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

var openaiAPIKey = mustGetenv("OPENAI_API_KEY")
var ImageFileBase = mustGetenv("IMAGE_FILE_BASE")
var rootURL = mustGetenv("ROOT_URL")

type GeneratedImage struct {
	model.GeneratedImage
}

func (gi *GeneratedImage) URL() string {
	return fmt.Sprintf("%s/generated_images/%d", rootURL, gi.ID)
}

func GenerateDALLE(ctx context.Context, prompt string) (*GeneratedImage, error) {
	client := openai.NewClient(openaiAPIKey)

	req := openai.ImageRequest{
		Prompt:         prompt,
		Model:          openai.CreateImageModelDallE3,
		N:              1,
		Size:           "1024x1024",
		ResponseFormat: "url",
	}

	imgResp, err := client.CreateImage(ctx, req)
	if err != nil {
		return nil, err
	}

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := model.New(tx)
	gi, err := q.CreateGeneratedImage(ctx, model.CreateGeneratedImageParams{
		Prompt:        prompt,
		RevisedPrompt: imgResp.Data[0].RevisedPrompt,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(imgResp.Data[0].URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = os.MkdirAll(ImageFileBase, os.FileMode(0755))
	if err != nil {
		return nil, err
	}

	imgFile, err := os.Create(fmt.Sprintf("%s/%d.png", ImageFileBase, gi.ID))
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	_, err = io.Copy(imgFile, resp.Body)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &GeneratedImage{gi}, nil
}
