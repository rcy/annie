package epigram

import (
	_ "embed"
	"goirc/internal/responder"
	"math/rand"
	"strings"
)

var (
	//go:embed data/epigrams.txt
	epigramsContent string
	epigrams        []string = strings.Split(epigramsContent, "\n")
)

func Handle(params responder.Responder) error {
	ri := rand.Intn(len(epigrams))

	params.Privmsgf(params.Target(), "%s", epigrams[ri])

	return nil
}
