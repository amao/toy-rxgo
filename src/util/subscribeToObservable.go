package util

import "github.com/amao/toy-rxgo/src/base"

func SubscribeToObservable(obs base.Subscribable) func(base.Observer) base.Unsubscribable {
	return func(subscriber base.Observer) base.Unsubscribable {
		return obs.Subscribe(subscriber.Next, subscriber.Error, subscriber.Complete)
	}
}
