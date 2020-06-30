package observables

import (
	"time"

	"github.com/amao/toy-rxgo/src/base"
)

func Interval(period int) base.Observable {
	return base.NewObservable(func(observer base.Observer) base.Unsubscribable {
		i := 0
		observer.Next(i)
		i++
		t := time.NewTicker(time.Duration(period) * time.Millisecond)
		go func() {
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
