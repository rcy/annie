package idle

import (
	"time"
)

// Run handler repeatedly after timeout when idle.
// Returns a function which when called resets the timeout.
func Repeat(timeout time.Duration, handler func()) func() {
	ch := make(chan bool)

	go func() {
		var alarm time.Time
		reset := func() {
			alarm = time.Now().Add(timeout)
		}
		reset()

		for {
			select {
			case <-ch:
				reset()
			default:
				time.Sleep(time.Second)
			}

			if time.Now().After(alarm) {
				handler()
				reset()
			}
		}
	}()

	return func() {
		ch <- true
	}
}

// Run handler after timeout when idle, and repeat only after reset is called.
// Returns a function which when called resets the timeout.
func RepeatAfterReset(timeout time.Duration, handler func()) func() {
	ch := make(chan bool)

	go func() {
		var alarm time.Time
		var canRun bool

		reset := func() {
			alarm = time.Now().Add(timeout)
			canRun = true
		}
		reset()

		for {
			select {
			case <-ch:
				reset()
			default:
				time.Sleep(time.Second)
			}

			if canRun && time.Now().After(alarm) {
				handler()
				canRun = false
			}
		}
	}()

	return func() {
		ch <- true
	}
}
