package bot

import (
	"goirc/internal/responder"
	"goirc/pubsub"
	"strings"

	irc "github.com/thoj/go-ircevent"
)

type HandlerParams struct {
	privmsgf  func(string, string, ...interface{})
	msg       string
	nick      string
	target    string
	matches   []string
	LastEvent *irc.Event
}

func NewHandlerParams(target string, privmsgf func(string, string, ...interface{})) HandlerParams {
	return HandlerParams{target: target, privmsgf: privmsgf}
}

func (hp HandlerParams) Privmsgf(target string, format string, a ...interface{}) {
	hp.privmsgf(target, format, a...)
}

func (hp HandlerParams) Target() string {
	return hp.target
}

func (hp HandlerParams) Nick() string {
	return hp.nick
}

func (hp HandlerParams) Match(num int) string {
	return strings.TrimSpace(hp.matches[num])
}

func (hp HandlerParams) Matches() []string {
	return hp.matches
}

func (hp HandlerParams) Msg() string {
	return hp.msg
}

func (hp *HandlerParams) Publish(eventName string, payload any) {
	pubsub.Publish(eventName, payload)
}

type HandlerFunction func(responder.Responder) error
