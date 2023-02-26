package handlers

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func link(params Params) bool {
	re := regexp.MustCompile(`(https?://\S+)`)
	matches := re.FindSubmatch([]byte(params.Msg))

	if len(matches) > 0 {
		if params.Target == params.Nick {
			params.Privmsgf(params.Target, "not your personal secretary")
			return false
		}

		url := string(matches[1])

		err := insertNote(params.Db, params.Target, params.Nick, "link", url)
		if err != nil {
			log.Print(err)
			params.Privmsgf(params.Target, err.Error())
		} else {
			log.Printf("recorded url %s", url)
		}

		// post to twitter
		nvurl := os.Getenv("NICHE_VOMIT_URL")
		if nvurl != "" {
			res, err := http.Post(nvurl, "text/plain", strings.NewReader(url))
			if res.StatusCode >= 300 || err != nil {
				log.Printf("error posting to twitter %d %v\n", res.StatusCode, err)
			}
		}

		return true
	}
	return false
}
