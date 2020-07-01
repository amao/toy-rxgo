package operators

import (
	"math"

	"github.com/amao/toy-rxgo/src/base"
)

type delayState struct {
	source      delaySubscriber
	destination base.Subscriber
	scheduler   base.TimestampProviderAndSchedulerLike
}

func newDelayState(source delaySubscriber, destination base.Subscriber, scheduler base.TimestampProviderAndSchedulerLike) delayState {
	newInstance := delayState{}
	newInstance.source = source
	newInstance.destination = destination
	newInstance.scheduler = scheduler
	return newInstance
}

type delayMessage struct {
	time  float64
	value interface{}
}

func newDelayMessage(time float64, value interface{}) delayMessage {
	newInstance := delayMessage{}
	newInstance.time = time
	newInstance.value = value
	return newInstance
}

type delaySubscriber struct {
	parent    base.Subscriber
	queue     []delayMessage
	active    bool
	delay     float64
	scheduler base.TimestampProviderAndSchedulerLike
}

func newDelaySubscriber(destination base.Subscriber, delay float64, scheduler base.TimestampProviderAndSchedulerLike) delaySubscriber {
	newInstance := delaySubscriber{}
	newInstance.parent = base.NewSubscriber(destination.Next, destination.Error, destination.Complete)
	return newInstance
}

func dispatch(schedulerAction base.SchedulerAction, state interface{}) {
	ds := state.(delayState)
	source := ds.source
	queue := source.queue
	scheduler := ds.scheduler
	destination := ds.destination

	for len(queue) > 0 && (queue[0].time-scheduler.Now() <= 0) {
		destination.Next(queue[0].value)
		queue = queue[1:]
	}

	if len(queue) > 0 {
		delay := math.Max(0, queue[0].time-scheduler.Now())
		schedulerAction.Schedule(state, delay)
	} else if source.parent.IsStopped {
		source.parent.Complete()
		source.active = false
	} else {
		schedulerAction.Unsubscribe()
		source.active = false
	}
}

func (d delaySubscriber) schedule(scheduler base.TimestampProviderAndSchedulerLike) {
	d.active = true
	destination := d.parent.Subscription
	destination.Add(
		scheduler.Schedule(dispatch, d.delay, newDelayState(d, d.parent, scheduler)),
	)
}

func (d delaySubscriber) Next(value interface{}) {
	scheduler := d.scheduler
	message := newDelayMessage(scheduler.Now()+d.delay, value)
	d.queue = append(d.queue, message)
	if !d.active {
		d.schedule(scheduler)
	}
}

func (d delaySubscriber) Error(err error) {
	if len(d.queue) == 0 {
		d.parent.Destination.Complete()
	}

	d.parent.Unsubscribe()
}

func (d delaySubscriber) Complete() {
	d.parent.Unsubscribe()
}

type delayOperator struct {
	delay     float64
	scheduler base.TimestampProviderAndSchedulerLike
}

func newDelayOperator(delay float64, scheduler base.TimestampProviderAndSchedulerLike) delayOperator {
	newInstance := delayOperator{}
	newInstance.delay = delay
	newInstance.scheduler = scheduler
	return newInstance
}

func (d delayOperator) Call(subscriber base.Subscriber, source base.Observable) base.Unsubscribable {
	do := newDelaySubscriber(subscriber, d.delay, d.scheduler)
	return source.Subscribe(do.Next, do.Error, do.Complete)
}

//todo
