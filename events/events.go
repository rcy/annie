package events

type BotJoined struct {
	Nick string
}

func (BotJoined) EventType() string { return "BotJoined" }
func (BotJoined) Aggregate() string { return "channel" }

type BotParted struct {
	Nick string
}

func (BotParted) EventType() string { return "BotParted" }
func (BotParted) Aggregate() string { return "channel" }

type BotQuit struct {
	Nick string
}

func (BotQuit) EventType() string { return "BotQuit" }
func (BotQuit) Aggregate() string { return "channel" }

type NickJoined struct {
	Nick string
}

func (NickJoined) EventType() string { return "NickJoined" }
func (NickJoined) Aggregate() string { return "channel" }

type NickParted struct {
	Nick string
}

func (NickParted) EventType() string { return "NickParted" }
func (NickParted) Aggregate() string { return "channel" }

type NickQuit struct {
	Nick string
}

func (NickQuit) EventType() string { return "NickQuit" }
func (NickQuit) Aggregate() string { return "channel" }

type NamesListed struct {
	Nicks []string
}

func (NamesListed) EventType() string { return "NamesListed" }
func (NamesListed) Aggregate() string { return "channel" }

type PublicMessageReceived struct {
	Nick    string
	Content string
}

func (PublicMessageReceived) EventType() string { return "PublicMessageReceived" }
func (PublicMessageReceived) Aggregate() string { return "channel" }

type PrivateMessageReceived struct {
	Content string
}

func (PrivateMessageReceived) EventType() string { return "PrivateMessageReceived" }
func (PrivateMessageReceived) Aggregate() string { return "nick" }
