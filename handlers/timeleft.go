package handlers

import (
	"goirc/internal/responder"
	"math"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TimeLeft(params responder.Responder) error {
	left := time.Until(time.Unix(2<<30, 0))

	years := int(math.Round(left.Hours() / 24 / 365))
	months := int(math.Round(left.Hours() / 24 / 365 * 12))
	weeks := int(math.Round(left.Hours() / 24 / 7))
	days := int(math.Round(left.Hours() / 24))
	hours := int(math.Round(left.Hours()))
	minutes := int(math.Round(left.Minutes()))
	seconds := int(math.Round(left.Seconds()))

	p := message.NewPrinter(language.English)

	str := p.Sprintf("%d years / %d months / %d weeks / %d days / %d hours / %d minutes / %d seconds\n", years, months, weeks, days, hours, minutes, seconds)

	params.Privmsgf(params.Target(), "%s", str)

	return nil
}
