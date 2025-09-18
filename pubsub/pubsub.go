package pubsub

type subscription struct {
	name    string
	handler func(any)
}

var subs []subscription

func Subscribe(name string, handler func(payload any)) {
	subs = append(subs, subscription{
		name:    name,
		handler: handler,
	})
}

func Publish(name string, payload any) {
	for _, sub := range subs {
		if sub.name == name {
			sub.handler(payload)
		}
	}
}
