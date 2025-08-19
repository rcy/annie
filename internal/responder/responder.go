package responder

type Responder interface {
	Privmsgf(string, string, ...interface{})
	Target() string
	Nick() string
	Match(num int) string
	Matches() []string
	Msg() string
}
