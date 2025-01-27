package bible

import (
	"errors"
	"fmt"
	"goirc/fetch"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/solafide-dev/gobible"
	"github.com/solafide-dev/gobible/bible"
)

type Bible struct {
	gobible           *gobible.GoBible
	activeTranslation *bible.Bible
}

func New() Bible {
	return Bible{
		gobible: gobible.NewGoBible(),
	}
}

func (b *Bible) SetActiveTranslation(translation string) error {
	_, bytes, err := fetch.Get(fmt.Sprintf("https://raw.githubusercontent.com/solafide-dev/gobible-gen/refs/heads/master/generated/%s.json", translation), 1000*time.Hour)
	if err != nil {
		return err
	}
	err = b.gobible.LoadString(string(bytes))
	if err != nil {
		return err
	}

	bib, err := b.gobible.GetTranslation(translation)
	if err != nil {
		return err
	}

	b.activeTranslation = bib

	return nil
}

type ref struct {
	book       string
	chapter    int
	startVerse int
	endVerse   int
}

var refRegexp = regexp.MustCompile("^(.+) ([0-9]+):([0-9]+)(-([0-9]+))?")

func parseRef(r string) (*ref, error) {
	matches := refRegexp.FindStringSubmatch(r)
	if matches == nil {
		return nil, fmt.Errorf("could not parse reference")
	}

	book := matches[1]
	chapter, _ := strconv.Atoi(matches[2])
	startVerse, _ := strconv.Atoi(matches[3])
	endVerse, _ := strconv.Atoi(matches[5])
	if endVerse == 0 {
		endVerse = startVerse
	}
	if startVerse > endVerse {
		return nil, errors.New("startVerse is greater than endVerse")
	}

	return &ref{
		book:       book,
		chapter:    chapter,
		startVerse: startVerse,
		endVerse:   endVerse,
	}, nil
}

func (b Bible) Lookup(r string) (string, error) {
	ref, err := parseRef(r)
	if err != nil {
		return "", err
	}

	start := float64(ref.startVerse)
	end := float64(ref.endVerse)
	verses := make([]string, 1+int(math.Max(start, end)-math.Min(start, end)))

	for i := range verses {
		verse, err := b.GetVerse(ref.book, ref.chapter, int(start)+i)
		if err != nil {
			return "", err
		}
		verses[i] = verse
	}
	return strings.Join(verses, " "), nil
}

func (b Bible) GetVerse(book string, chapter int, verse int) (string, error) {
	if b.activeTranslation == nil {
		return "", errors.New("no active translation set")
	}

	bk := b.activeTranslation.GetBook(book)
	if bk == nil {
		return "", errors.New("book not found")
	}
	ch := bk.GetChapter(chapter)
	if ch == nil {
		return "", errors.New("chapter not found")
	}
	vs := ch.GetVerse(verse)
	if vs == nil {
		return "", errors.New("verse not found")
	}

	return vs.Text, nil
}
