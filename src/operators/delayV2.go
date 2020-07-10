package operators

import (
	"time"

	"github.com/amao/toy-rxgo/src/base"
)

type delayState struct {
	source      *delaySubscriberV2
	destination base.SubscriberLike
	scheduler   base.SchedulerLike
}

type delayMessage struct {
	time  time.Time
	value interface{}
}

func newDelayMessage(time time.Time, value interface{}) delayMessage {
	newInstance := new(delayMessage)
	newInstance.time = time
	newInstance.value = value
	return *newInstance
}

type delaySubscriberV2 struct {
	*base.Subscriber
	queue     []delayMessage
	active    bool //default value =false
	delay     uint
	scheduler base.SchedulerLike
}

func dispatch(schedulerAction base.SchedulerAction, state interface{}) {
	delayState := state.(delayState)
	source := delayState.source
	queue := source.queue
	scheduler := delayState.scheduler
	destination := delayState.destination

	for len(queue) > 0 && queue[0].time.Sub(scheduler.Now()) <= 0 {
		destination.Next(queue[0].value)
		queue = queue[1:]
	}

	if len(queue) > 0 {
		diff := queue[0].time.Sub(scheduler.Now()).Milliseconds()
		var delay uint
		if diff <= 0 {
			// must be greater than zero
			delay = 1
		} else {
			delay = uint(diff)
		}
		schedulerAction.Schedule(state, delay)
	} else if source.IsStopped {
		source.Destination.Complete()
		source.active = false
	} else {
		schedulerAction.Unsubscribe()
		source.active = false
	}
}

func newDelaySubscriberV2(destination base.SubscriberLike, delay uint, scheduler base.SchedulerLike) delaySubscriberV2 {
	newInstance := new(delaySubscriberV2)
	parentSubscriber := base.NewSubscriber(destination)
	newInstance.Subscriber = &parentSubscriber
	newInstance.delay = delay
	newInstance.scheduler = scheduler
	return *newInstance
}

func (d *delaySubscriberV2) _schedule(scheduler base.SchedulerLike) {
	d.active = true
	destination := d.Destination
	subscriptionLike := scheduler.Schedule(scheduler, dispatch, d.delay, delayState{
		source:      d,
		destination: destination,
		scheduler:   scheduler,
	})
	destination.Add(subscriptionLike)
}

func (d *delaySubscriberV2) Next(value interface{}) {
	if !d.IsStopped {
		scheduler := d.scheduler
		messge := newDelayMessage(scheduler.Now().Add(time.Millisecond*time.Duration(d.delay)), value)
		d.queue = append(d.queue, messge)
		if d.active == false {
			d._schedule(scheduler)
		}
	}
}

func (d *delaySubscriberV2) Error(err error) {
	d.queue = nil
	d.Destination.Error(err)
	d.Unsubscribe()
}

func (d *delaySubscriberV2) Complete() {
	if len(d.queue) == 0 {
		d.Destination.Complete()
	}

	d.Unsubscribe()
}

type delayOperatorV2 struct {
	delay     uint
	scheduler base.SchedulerLike
}

func newDelayOperatorV2(delay uint, scheduler base.SchedulerLike) delayOperatorV2 {
	newInstance := new(delayOperatorV2)
	newInstance.delay = delay
	newInstance.scheduler = scheduler
	return *newInstance
}

func (d *delayOperatorV2) Call(subscriber base.SubscriberLike, source base.Observable) base.SubscriptionLike {
	ndsv2 := newDelaySubscriberV2(subscriber, d.delay, d.scheduler)
	return source.Subscribe(&ndsv2)
}

// do not use v2
func delayV2(delay uint, args ...interface{}) base.OperatorFunction {
	result := func(source base.Observable) base.Observable {
		async := base.NewAsyncScheduler("AsyncAction", func() time.Time {
			return time.Now()
		})
		op := newDelayOperatorV2(delay, &async)
		return source.Lift(&op)
	}
	return result
}
