package operators

import (
	"sync/atomic"
	"time"

	"github.com/amao/toy-rxgo/src/base"
)

type message struct {
	time  time.Time
	value interface{}
}

type delaySubscriber struct {
	super *base.Subscriber
	*base.Subscriber
	delay float64
	queue []message
}

func newDelaySubscriber(destination base.SubscriberLike, delay float64) delaySubscriber {
	newInstance := new(delaySubscriber)
	super := base.NewSubscriber(destination)
	self := base.NewSubscriber(destination)
	newInstance.Subscriber = &self
	newInstance.super = &super
	newInstance.Destination.Add(newInstance)
	newInstance.delay = delay
	var jobIsRunning uint32
	newInstance.SetInnerNext(func(value interface{}) {
		newInstance.queue = append(newInstance.queue, message{
			time:  time.Now(),
			value: value,
		})
		if atomic.CompareAndSwapUint32(&jobIsRunning, 0, 1) {
			go func() {
				for {
					for len(newInstance.queue) > 0 && float64(time.Now().UnixNano()/1e6-newInstance.queue[0].time.UnixNano()/1e6) >= newInstance.delay {
						newInstance.Destination.Next(newInstance.queue[0].value)
						newInstance.queue = newInstance.queue[1:]
					}
					if len(newInstance.queue) > 0 {
						continue
					}
					break
				}
				newInstance.Complete()
				atomic.StoreUint32(&jobIsRunning, 0)
			}()
		}
	})
	newInstance.SetInnerError(func(err error) {
		newInstance.Destination.Error(err)
		newInstance.Unsubscribe()
	})
	newInstance.SetInnerComplete(func() {
		go func() {
			time.Sleep(time.Duration(newInstance.delay) * time.Millisecond)
			newInstance.Destination.Complete()
			newInstance.Unsubscribe()
		}()
	})
	return *newInstance
}

type delayOperator struct {
	delay float64
}

func newDelayOperator(delay float64) delayOperator {
	newInstance := new(delayOperator)
	newInstance.delay = delay
	return *newInstance
}

func (d *delayOperator) Call(subscriber base.SubscriberLike, source base.Subscribable) base.SubscriptionLike {
	nds := newDelaySubscriber(subscriber, d.delay)
	return source.Subscribe(nds)
}

func Delay(delay float64, args ...interface{}) base.OperatorFunction {
	result := func(source base.Observable) base.Observable {
		op := newDelayOperator(delay)
		return source.Lift(&op)
	}
	return result
}
