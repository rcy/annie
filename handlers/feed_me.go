package handlers

import (
	"database/sql"
	"fmt"
	"goirc/bot"
	"goirc/durfmt"
	"goirc/internal/idstr"
	"goirc/model"
	"goirc/model/notes"
	"goirc/util"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/kljensen/snowball"
	"github.com/kljensen/snowball/english"
)

func candidateLinks(age time.Duration) ([]notes.Note, error) {
	notes := []notes.Note{}

	query := `select id, created_at, nick, text, kind from notes where created_at <= datetime(?) and nick = target order by random() limit 420`

	t := time.Now().UTC().Add(-age).Format(time.RFC3339)
	err := model.DB.Select(&notes, query, t)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return notes, nil
}

const (
	MINAGE    = time.Hour * 24
	THRESHOLD = 5
	COOLOFF   = time.Hour * 5
)

var lastSentAt = time.Unix(0, 0)

func canSendIn(startTime time.Time) time.Duration {
	return startTime.Add(COOLOFF).Sub(time.Now())
}

func FeedMe(params bot.HandlerParams) error {
	waitFor := canSendIn(lastSentAt)

	if waitFor > 0 {
		if params.Nick != "" {
			params.Privmsgf(params.Target, "throttled for another %s", durfmt.Format(waitFor))
		}
		return nil
	}

	target := ""
	if params.LastEvent != nil {
		target = params.LastEvent.Message()
	}

	notes, err := candidateLinks(MINAGE)
	if err != nil {
		return err
	}

	if len(notes) < THRESHOLD {
		if params.Nick != "" {
			params.Privmsgf(params.Target, "not enough links to feed the channel")
		}
		return nil
	}

	candidates := make([]string, len(notes))
	for i, n := range notes {
		candidates[i] = n.Text
	}
	bestIndex := bestMatch(target, candidates)

	note := notes[bestIndex]

	_, err = model.DB.Exec(`update notes set target = ? where id = ?`, params.Target, note.Id)
	if err != nil {
		return err
	}

	ready, fermenting, err := health()
	if err != nil {
		return err
	}

	var text string
	if note.Kind == "link" {
		str, err := idstr.Encode(note.Id)
		if err != nil {
			return err
		}
		text = fmt.Sprintf("%s/%s", os.Getenv("ROOT_URL"), str)
	} else {
		text = note.Text
	}

	params.Privmsgf(params.Target, "%s (%s ago) [pipe=%d+%d]", text, util.Since(note.CreatedAt), ready, fermenting)

	lastSentAt = time.Now()

	return nil
}

func PipeHealth(params bot.HandlerParams) error {
	ready, fermenting, err := health()
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%d links ready to serve (%d fermenting)", ready, fermenting)

	return nil
}

// Return ready, fermenting, error
func health() (int, int, error) {
	readyNotes, err := candidateLinks(MINAGE)
	if err != nil {
		return 0, 0, err
	}
	totalNotes, err := candidateLinks(0)
	if err != nil {
		return 0, 0, err
	}

	ready := len(readyNotes)
	fermenting := len(totalNotes) - ready

	return ready, fermenting, nil
}

// Return the index of the match between s and candidates
func bestMatch(s string, candidates []string) int {
	target := stemMessage(cleanMessage(s))

	type cand struct {
		index int
		score int
	}

	cas := []cand{}

	for i, c := range candidates {
		stems := stemMessage(cleanMessage(c))
		score := compareArray(target, stems)
		cas = append(cas, cand{i, score})
	}

	sort.Slice(cas, func(i, j int) bool {
		return cas[i].score > cas[j].score
	})

	return cas[0].index
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func cleanMessage(msg string) string {
	return nonAlphanumericRegex.ReplaceAllString(msg, " ")
}

func stemMessage(msg string) []string {
	stems := []string{}

	for _, word := range strings.Fields(msg) {
		if !english.IsStopWord(word) {
			stem, _ := snowball.Stem(word, "english", false)
			stems = append(stems, stem)
		}
	}

	return stems
}

func compareArray(arr1, arr2 []string) int {
	matches := 0

outer:
	for _, e1 := range arr1 {
		for _, e2 := range arr2 {
			if e1 == e2 {
				matches += 1
				continue outer
			}
		}
	}

	return matches
}
