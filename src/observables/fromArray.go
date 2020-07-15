package observables

import "github.com/amao/toy-rxgo/src/base"

func FromArray(array []interface{}) base.Observable {
	ob := base.NewObservable(func(observer base.SubscriberLike) base.SubscriptionLike {
		for _, value := range array {
			observer.Next(value)
		}
		observer.Complete()
		sp := base.NewSubscription(nil)
		return &sp
	})
	return ob
}
