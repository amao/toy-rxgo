package observables

import (
	"time"

	"github.com/amao/toy-rxgo/src/base"
)

func Timer(duration int, period interface{}) *base.Observable {
	result := base.NewObservable(func(observer base.SubscriberLike) base.SubscriptionLike {
		if period == nil {
			time.Sleep(time.Duration(duration) * time.Millisecond)
			observer.Next(0)
			observer.Complete()
			sp := base.NewSubscription(nil)
			return &sp
		} else {
			t := time.NewTicker(time.Duration(period.(int)) * time.Millisecond)
			go func() {
				i := 0
				time.Sleep(time.Duration(duration) * time.Millisecond)
				observer.Next(i)
				i++
				for range t.C {
					observer.Next(i)
					i++
				}
			}()

			sp := base.NewSubscription(func() {
				t.Stop()
			})

			return &sp
		}
	})

	return &result
}
