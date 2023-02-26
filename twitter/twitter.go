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
		if res.StatusCode >= 300 || err != nil {
			return errors.New(fmt.Sprintf("error posting to twitter %d %v\n", res.StatusCode, err))
		}
	}
	return nil
}
