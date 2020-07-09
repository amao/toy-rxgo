package observables

import "github.com/amao/toy-rxgo/src/base"

func FromArray(array []int) base.Observable {
	return base.NewObservable(func(observer base.SubscriberLike) base.SubscriptionLike {
		for _, value := range array {
			observer.Next(value)
		}
		observer.Complete()
		sp := base.NewSubscription(nil)
		return &sp
	})
}
