package twitter

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Post(text string) error {
	nvurl := os.Getenv("NICHE_VOMIT_URL")
	if nvurl != "" {
		res, err := http.Post(nvurl, "text/plain", strings.NewReader(text))
		if err != nil {
			return errors.New(fmt.Sprintf("error posting to twitter err=%s\n", err))
		}
		if res.StatusCode >= 300 {
			return errors.New(fmt.Sprintf("error posting to twitter statusCode=%s\n", err))
		}
	}
	return nil
}
