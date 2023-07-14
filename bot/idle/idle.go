package idle

import (
	"time"
)

var lastMessage time.Time

func Reset() {
	lastMessage = time.Now()
}

func Every(duration time.Duration, fn func()) {
	Reset()
	for {
		time.Sleep(1 * time.Minute)

		if time.Now().Sub(lastMessage) >= duration {
			Reset()
			fn()
		}
	}
}
