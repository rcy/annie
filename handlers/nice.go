package handlers

import (
	"goirc/internal/responder"
	"time"
)

func Nice(params responder.Responder) error {
	go func() {
		time.Sleep(10 * time.Second)
		params.Privmsgf(params.Target(), "%s: nice", params.Nick())
	}()

	return nil
}
