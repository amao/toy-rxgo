package observables

import "github.com/amao/toy-rxgo/src/base"

func FromArray(array []int) base.Observable {
	return base.NewObservable(func(observer base.Observer) base.Unsubscribable {
		for _, value := range array {
			observer.Next(value)
		}
		observer.Complete()
		return base.NewSubscription()
	})
}
