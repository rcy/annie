package idle

import (
	"time"
)

var lastMessage struct {
	sentAt time.Time
	sentBy string
}

func Reset(nick string) {
	lastMessage.sentAt = time.Now()
	lastMessage.sentBy = nick
}

func Every(duration time.Duration, fn func()) {
	Reset("nobody")
	for {
		time.Sleep(1 * time.Minute)

		if time.Since(lastMessage.sentAt) >= duration {
			fn()
		}
	}
}
