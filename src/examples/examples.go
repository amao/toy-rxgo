package main

import (
	"fmt"
	"time"

	"github.com/amao/toy-rxgo/src/base"
	"github.com/amao/toy-rxgo/src/observables"
	"github.com/amao/toy-rxgo/src/operators"
)

func main() {
	subscriber := base.NewSubscriber(
		func(value interface{}) {
			fmt.Println(value)
		},
	)

	subscription := observables.Interval(1000).Pipe(
		operators.Map(func(x interface{}) interface{} { return x.(int) * 10 }),
		operators.Filter(func(x interface{}) bool { return x.(int) < 60 }),
		operators.SwitchMap(func(x interface{}) base.Subscribable {
			return observables.Of(fmt.Sprintf("%d switchmap to %d", x, x.(int)*10))
		}),
	).Subscribe(
		subscriber.Next,
		subscriber.Error,
		subscriber.Complete,
	)

	for {
		time.Sleep(10 * time.Second)
		break
	}

	subscription.Unsubscribe()

	for {
	}
}
