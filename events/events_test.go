package events

import (
	"testing"
)

func TestPublishSubscribe(t *testing.T) {
	var msgHandled = false
	var msg2Handled = false
	var otherHandled = false

	Subscribe("msg", func(payload any) {
		msgHandled = true
	})

	Subscribe("msg", func(payload any) {
		msg2Handled = true
	})

	Subscribe("other", func(payload any) {
		otherHandled = true
	})

	Publish("msg", true)

	if !msgHandled {
		t.Error("msg not true")
	}
	if !msg2Handled {
		t.Error("msg2 not true")
	}
	if otherHandled {
		t.Error("otherHandled not false")
	}
}
