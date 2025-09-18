package events

type channelEvent struct{}

func (channelEvent) Aggregate() string { return "channel" }

type BotJoined struct {
	channelEvent
	Nick string
}

type BotParted struct {
	channelEvent
	Nick string
}

type BotQuit struct {
	channelEvent
	Nick string
}

type NickJoined struct {
	channelEvent
	Nick string
}

type NickParted struct {
	channelEvent
	Nick string
}

type NickQuit struct {
	channelEvent
	Nick string
}

type NamesListed struct {
	channelEvent
	Nicks []string
}

type PublicMessageReceived struct {
	channelEvent
	Nick    string
	Content string
}

type nickEvent struct{}

func (nickEvent) Aggregate() string { return "nick" }

type PrivateMessageReceived struct {
	nickEvent
	Content string
}
