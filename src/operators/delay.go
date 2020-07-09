package operators

import (
	"time"

	"github.com/amao/toy-rxgo/src/base"
)

type delaySubscriber struct {
	*base.Subscriber
	delay float64
}

func newDelaySubscriber(destination base.SubscriberLike, delay float64) delaySubscriber {
	newInstance := new(delaySubscriber)
	s := base.NewSubscriber(destination.Next, destination.Error, destination.Complete)
	newInstance.Subscriber = &s
	newInstance.delay = delay
	return *newInstance
}

func (d *delaySubscriber) Next(value interface{}) {
	d.Destination.Next(value)
}

func (d *delaySubscriber) Error(err error) {
	d.Destination.Error(err)
	d.Subscriber.Unsubscribe()
}

func (d *delaySubscriber) Complete() {
	d.Destination.Complete()
	d.Subscriber.Unsubscribe()
}

type delayOperator struct {
	delay float64
}

func newDelayOperator(delay float64) delayOperator {
	newInstance := new(delayOperator)
	newInstance.delay = delay
	return *newInstance
}

func (d *delayOperator) Call(subscriber base.SubscriberLike, source base.Observable) base.SubscriptionLike {
	nds := newDelaySubscriber(subscriber, d.delay)
	time.Sleep(time.Duration(d.delay) * time.Millisecond)
	return source.Subscribe(nds.Next, nds.Error, nds.Complete)
}

func Delay(delay float64, args ...interface{}) base.OperatorFunction {
	result := func(source base.Observable) base.Observable {
		op := newDelayOperator(delay)
		return source.Lift(&op)
	}
	return result
}
