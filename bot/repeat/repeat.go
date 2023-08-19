package repeat

import (
	"time"
)

func Every(duration time.Duration, fn func()) {
	for {
		time.Sleep(duration)

		fn()
	}
}
