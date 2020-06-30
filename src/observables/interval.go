package observables

import (
	"time"

	"github.com/amao/toy-rxgo/src/base"
)

func Interval(period int) base.Observable {
	return base.NewObservable(func(observer base.Observer) base.Unsubscribable {

		t := time.NewTicker(time.Duration(period) * time.Millisecond)
		go func() {
			i := 0
			for range t.C {
				observer.Next(i)
				i++
			}
		}()

		return base.NewSubscription(func() {
			t.Stop()
		})
	})
}
