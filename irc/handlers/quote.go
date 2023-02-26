package handlers

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func quote(params Params) bool {
	// match anything that starts with a quote and has no subsequent quotes
	re := regexp.MustCompile(`^("[^"]+)$`)
	matches := re.FindSubmatch([]byte(params.Msg))

	if len(matches) > 0 {
		if params.Target == params.Nick {
			params.Privmsgf(params.Target, "not your personal secretary")
			return false
		}

		text := string(matches[1])

		err := insertNote(params.Db, params.Target, params.Nick, "quote", text)
		if err != nil {
			log.Print(err)
		}

		// post to twitter
		nvurl := os.Getenv("NICHE_VOMIT_URL")
		if nvurl != "" {
			res, err := http.Post(nvurl, "text/plain", strings.NewReader(text))
			if res.StatusCode >= 300 || err != nil {
				log.Printf("error posting to twitter %d %v\n", res.StatusCode, err)
			}
		}

		return true
	}
	return false
}
