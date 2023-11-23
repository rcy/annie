package bot

import irc "github.com/thoj/go-ircevent"

type HandlerParams struct {
	Privmsgf  func(string, string, ...interface{})
	Msg       string
	Nick      string
	Target    string
	Matches   []string
	LastEvent *irc.Event
}

type HandlerFunction func(HandlerParams) error
