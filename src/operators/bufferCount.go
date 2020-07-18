package operators

import "github.com/amao/toy-rxgo/src/base"

type bufferCountSubscriber struct {
	*base.Subscriber
	super      *base.Subscriber
	buffer     []interface{}
	bufferSize int
}

func newBufferCountSubscriber(destination base.SubscriberLike, bufferSize int) bufferCountSubscriber {
	newInstance := new(bufferCountSubscriber)
	self := base.NewSubscriber(destination)
	super := base.NewSubscriber(destination)
	newInstance.Subscriber = &self
	newInstance.super = &super
	newInstance.Destination.Add(newInstance)
	newInstance.bufferSize = bufferSize

	newInstance.SetInnerNext(func(value interface{}) {

		newInstance.buffer = append(newInstance.buffer, value)

		if len(newInstance.buffer) == newInstance.bufferSize {
			newInstance.Destination.Next(newInstance.buffer)
			newInstance.buffer = []interface{}{}
		}
	})

	newInstance.SetInnerComplete(func() {
		if len(newInstance.buffer) > 0 {
			newInstance.Destination.Next(newInstance.buffer)
		}
		super.CallInnerComplete()
	})

	return *newInstance
}

type bufferSkipCountSubscriber struct {
	*base.Subscriber
	super            *base.Subscriber
	buffers          [][]interface{}
	count            int
	bufferSize       int
	startBufferEvery int
}

func newBufferSkipCountSubscriber(destination base.SubscriberLike, bufferSize int, startBufferEvery int) bufferSkipCountSubscriber {
	newInstance := new(bufferSkipCountSubscriber)

	self := base.NewSubscriber(destination)
	super := base.NewSubscriber(destination)
	newInstance.Subscriber = &self
	newInstance.super = &super
	newInstance.Destination.Add(newInstance)
	newInstance.bufferSize = bufferSize
	newInstance.startBufferEvery = startBufferEvery

	newInstance.SetInnerNext(func(value interface{}) {

		newInstance.count++
		if newInstance.count%startBufferEvery == 0 {
			newInstance.buffers = append(newInstance.buffers, []interface{}{})
		}

		for i := len(newInstance.buffers) - 1; i >= 0; i-- {
			newInstance.buffers[i] = append(newInstance.buffers[i], value)
			if len(newInstance.buffers[i]) == newInstance.bufferSize {
				newInstance.Destination.Next(newInstance.buffers[i])
				newInstance.buffers = append(newInstance.buffers[:i], newInstance.buffers[i+1:]...)
			}
		}
	})

	newInstance.SetInnerComplete(func() {

		for len(newInstance.buffers) > 0 {
			buffer := newInstance.buffers[0]
			newInstance.buffers = newInstance.buffers[1:]
			if len(buffer) > 0 {
				newInstance.Destination.Next(buffer)
			}
		}

		newInstance.super.CallInnerComplete()
	})

	return *newInstance
}

type bufferCountOperator struct {
	subscriberClass  string
	bufferSize       int
	startBufferEvery interface{}
}

func newBufferCountOperator(bufferSize int, startBufferEvery interface{}) bufferCountOperator {
	newInstance := new(bufferCountOperator)
	newInstance.bufferSize = bufferSize

	if startBufferEvery == nil || bufferSize == startBufferEvery.(int) {
		newInstance.subscriberClass = "bufferCountSubscriber"
	} else {
		newInstance.subscriberClass = "bufferSkipCountSubscriber"
		newInstance.startBufferEvery = startBufferEvery
	}
	return *newInstance
}

func (b *bufferCountOperator) Call(subscriber base.SubscriberLike, source base.Subscribable) base.SubscriptionLike {
	var bs base.SubscriberLike

	if b.subscriberClass == "bufferCountSubscriber" {
		bs = newBufferCountSubscriber(subscriber, b.bufferSize)
	} else {
		bs = newBufferSkipCountSubscriber(subscriber, b.bufferSize, b.startBufferEvery.(int))
	}

	return source.Subscribe(bs)
}

func BufferCount(bufferSize int, startBufferEvery interface{}) base.OperatorFunction {
	result := func(source base.Observable) base.Observable {
		op := newBufferCountOperator(bufferSize, startBufferEvery)
		return source.Lift(&op)
	}
	return result
}
