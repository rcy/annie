package idle

import (
	"time"
)

// Run handler repeatedly after timeout given by duration.
// Returns a function which when called resets the timeout.
func Register(duration time.Duration, handler func()) func() {
	ch := make(chan bool)

	reset := func() {
		ch <- true
	}

	go func() {
		alarm := time.Now().Add(duration)
		for {
			select {
			case <-ch:
				alarm = time.Now().Add(duration)
			default:
				time.Sleep(time.Second)
				if time.Now().After(alarm) {
					handler()
					alarm = time.Now().Add(duration)
				}
			}
		}
	}()

	return reset
}
