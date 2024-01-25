package epigrams

import (
	_ "embed"
	"math/rand"
	"strings"
)

var (
	//go:embed data/epigrams.txt
	epigramsContent string
	epigrams        []string = strings.Split(epigramsContent, "\n")
)

func Random() string {
	ri := rand.Intn(len(epigrams))

	return epigrams[ri]
}
