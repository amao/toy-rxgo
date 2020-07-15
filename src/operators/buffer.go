package operators

import (
	"github.com/amao/toy-rxgo/src/base"
)

type bufferSubscriber struct {
	*base.OuterSubscriber
	buffer []interface{}
}

func newBufferSubscriber(destination base.SubscriberLike, closingNotifier base.Subscribable) bufferSubscriber {
	newInstance := new(bufferSubscriber)

	selfSubscriber := base.NewSubscriber(destination)
	selfOuterSubscriber := new(base.OuterSubscriber)
	newInstance.OuterSubscriber = selfOuterSubscriber
	newInstance.Subscriber = &selfSubscriber
	newInstance.Destination.Add(newInstance)

	newInstance.SetInnerNext(func(value interface{}) {
		newInstance.buffer = append(newInstance.buffer, value)
	})
	innerSubscriber := base.NewInnerSubscriber(newInstance, nil, nil)
	newInstance.Add(closingNotifier.Subscribe(innerSubscriber))

	return *newInstance
}

func (b *bufferSubscriber) NotifyNext(outerValue interface{}, innerValue interface{}, outerIndex interface{}, innerIndex interface{}, innerSub base.InnerSubscriber) {
	buffer := b.buffer
	b.buffer = nil
	b.Destination.Next(buffer)

}

type bufferOperator struct {
	closingNotifier *base.Observable
}

func newBufferOperator(closingNotifier base.Observable) bufferOperator {
	newInstance := new(bufferOperator)
	newInstance.closingNotifier = &closingNotifier
	return *newInstance
}

func (b *bufferOperator) Call(subscriber base.SubscriberLike, source base.Subscribable) base.SubscriptionLike {
	nbs := newBufferSubscriber(subscriber, b.closingNotifier)
	return source.Subscribe(nbs)
}

func Buffer(closingNotifier base.Observable) base.OperatorFunction {
	result := func(source base.Observable) base.Observable {
		op := newBufferOperator(closingNotifier)
		return source.Lift(&op)
	}
	return result
}
