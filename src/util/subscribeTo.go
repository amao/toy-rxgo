package util

import "github.com/amao/toy-rxgo/src/base"

func SubscribeTo(result base.Subscribable) func(base.Observer) base.Unsubscribable {
	return SubscribeToObservable(result)
}
