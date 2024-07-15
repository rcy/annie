package day

import (
	"context"
	"fmt"
	"testing"
)

func TestGenerateImage(t *testing.T) {
	ctx := context.Background()
	prompt := "A beautiful sunset over the mountains."
	url, err := generateImage(ctx, prompt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(url)
}
