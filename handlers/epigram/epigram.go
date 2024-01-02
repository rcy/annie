package epigram

import (
	_ "embed"
	"goirc/bot"
	"math/rand"
	"strings"
)

var (
	//go:embed data/epigrams.txt
	epigramsContent string
	epigrams        []string = strings.Split(epigramsContent, "\n")
)

func Handle(params bot.HandlerParams) error {
	ri := rand.Intn(len(epigrams))

	params.Privmsgf(params.Target, epigrams[ri])

	return nil
}
