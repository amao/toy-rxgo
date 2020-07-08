package observables

import (
	"time"

	"github.com/amao/toy-rxgo/src/base"
)

func Interval(period int) *base.Observable {
	result := base.NewObservable(func(observer base.Observer) base.SubscriptionLike {
		t := time.NewTicker(time.Duration(period) * time.Millisecond)
		go func() {
			i := 0
			for range t.C {
				observer.Next(i)
				i++
			}
		}()

		sp := base.NewSubscription(func() {
			t.Stop()
		})

		return &sp
	})

	return &result
}
