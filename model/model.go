package model

type ChannelNick struct {
	Channel   string
	Nick      string
	Present   string
	UpdatedAt string `db:"updated_at"`
}

type Later struct {
	RowId     int    `db:"rowid"`
	CreatedAt string `db:"created_at"`
	Nick      string
	Target    string
	Message   string
	Sent      bool
}
