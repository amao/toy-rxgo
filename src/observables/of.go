package observables

import "github.com/amao/toy-rxgo/src/base"

func Of(args ...interface{}) base.Observable {
	return base.NewObservable(func(observer base.Observer) base.Unsubscribable {
		for _, value := range args {
			observer.Next(value)
		}
		observer.Complete()
		return base.NewSubscription()
	})
}
