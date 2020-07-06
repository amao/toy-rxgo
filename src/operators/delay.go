package operators

import (
	"github.com/amao/toy-rxgo/src/base"
	"time"
)

type delayMessage struct {
	now   float64
	value interface{}
}

func newDelayMessage(value interface{}) delayMessage {
	newInstance := new(delayMessage)
	newInstance.now = float64(time.Now().Second() * 1000)
	newInstance.value = value
	return *newInstance
}

type delaySubscriber struct {
	*base.Subscriber
	queue []delayMessage
	delay float64
}

func newDelaySubscriber(destination base.Subscriber, delay float64) delaySubscriber {
	newInstance := new(delaySubscriber)
	s := base.NewSubscriber(destination.Next, destination.Error, destination.Complete)
	newInstance.Subscriber = &s
	newInstance.delay = delay

	go func() {
		for {
			for _, msg := range newInstance.queue {
				if float64(time.Now().Second()*1000)-msg.now >= delay {
					newInstance.Destination.Next(msg.value)
					newInstance.queue = newInstance.queue[1:]
				}
			}
		}
	}()

	return *newInstance
}

func (d *delaySubscriber) Next(value interface{}) {
	d.Destination.Next(value)
}

func (d *delaySubscriber) Error(err error) {
	d.queue = nil
	d.Destination.Error(err)
	d.Subscriber.Unsubscribe()
}

func (d *delaySubscriber) Complete() {
	if len(d.queue) == 0 {
		d.Destination.Complete()
	}
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

func (d *delayOperator) Call(subscriber *base.Subscriber, source base.Observable) base.Unsubscribable {
	nds := newDelaySubscriber(*subscriber, d.delay)
	time.Sleep(time.Duration(d.delay) * time.Millisecond)
	return source.Subscribe(nds.Next, nds.Error, nds.Complete)
}

func Delay(delay float64, args ...interface{}) base.OperatorFunction {
	result := func(source base.Observable) base.Observable {
		op := newDelayOperator(delay)
		return *source.Lift(&op)
	}
	return result
}
