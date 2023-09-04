package twitter

import (
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
			return fmt.Errorf("error posting to twitter %w\n", err)
		}
		if res.StatusCode >= 300 {
			return fmt.Errorf("error posting to twitter statusCode=%d\n", res.StatusCode)
		}
	}
	return nil
}
