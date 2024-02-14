package bot

import (
	"goirc/events"

	irc "github.com/thoj/go-ircevent"
)

type HandlerParams struct {
	Privmsgf  func(string, string, ...interface{})
	Msg       string
	Nick      string
	Target    string
	Matches   []string
	LastEvent *irc.Event
}

func (hp *HandlerParams) Publish(eventName string, payload any) {
	events.Publish(eventName, payload)
}

type HandlerFunction func(HandlerParams) error
