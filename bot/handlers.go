package bot

type HandlerParams struct {
	Privmsgf func(string, string, ...interface{})
	Msg      string
	Nick     string
	Target   string
	Matches  []string
}

type HandlerFunction func(HandlerParams) bool
