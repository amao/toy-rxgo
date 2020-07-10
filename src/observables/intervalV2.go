package observables

import (
	"time"

	"github.com/amao/toy-rxgo/src/base"
)

type IntervalState struct {
	subscriber base.SubscriberLike
	counter    uint
	period     uint
}

func dispatch(schedulerAction base.SchedulerAction, state interface{}) {
	intervalState := state.(IntervalState)
	subscriber := intervalState.subscriber
	counter := intervalState.counter
	period := intervalState.period
	subscriber.Next(counter)

	newIntervalState := IntervalState{
		subscriber: subscriber,
		counter:    counter + 1,
		period:     period,
	}

	schedulerAction.Schedule(newIntervalState, period)
}

func IntervalV2(period uint) *base.Observable {
	async := base.NewAsyncScheduler("AsyncAction", func() time.Time {
		return time.Now()
	})
	ob := base.NewObservable(func(subscriber base.SubscriberLike) base.SubscriptionLike {
		asyncAction := async.Schedule(&async, dispatch, period, IntervalState{
			subscriber: subscriber,
			counter:    0,
			period:     period,
		})
		subscriber.Add(asyncAction)
		return subscriber
	})
	return &ob
}
