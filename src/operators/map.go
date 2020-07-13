package operators

import (
	"github.com/amao/toy-rxgo/src/base"
)

func _map(transformFn func(interface{}) interface{}) func(base.Observable) base.Observable {
	result := func(inObservable base.Observable) base.Observable {
		outObservable := base.NewObservable(func(outObserver base.SubscriberLike) base.SubscriptionLike {
			inObserver := base.NewSubscriber(
				func(x interface{}) {
					y := transformFn(x)
					outObserver.Next(y)
				}, func(e error) {
					outObserver.Error(e)
				}, func() {})
			return inObservable.Subscribe(inObserver.Next, inObserver.Error, inObserver.Complete)
		})
		return outObservable
	}

	return result
}

type mapSubscriber struct {
	subscriber base.Subscriber
	project    func(interface{}) interface{}
}

func newMapSubscriber(destination base.SubscriberLike, project func(interface{}) interface{}) mapSubscriber {
	newInstance := mapSubscriber{}
	newInstance.subscriber = base.NewSubscriber(destination)
	newInstance.project = project
	return newInstance
}

func (s mapSubscriber) Next(value interface{}) {
	s.subscriber.Next(s.project(value))
}

func (s mapSubscriber) Error(err error) {
	s.subscriber.Error(err)
}

func (s mapSubscriber) Complete() {
	s.subscriber.Complete()
}

type mapOperator struct {
	project func(interface{}) interface{}
}

func newMapOperator(project func(interface{}) interface{}) mapOperator {
	newInstance := mapOperator{}
	newInstance.project = project
	return newInstance
}

func (m mapOperator) Call(subscriber base.SubscriberLike, source base.Subscribable) base.SubscriptionLike {
	nms := newMapSubscriber(subscriber, m.project)
	return source.Subscribe(nms.Next, nms.Error, nms.Complete)
}

func Map(project func(interface{}) interface{}) base.OperatorFunction {
	result := func(source base.Observable) base.Observable {
		op := newMapOperator(project)
		return source.Lift(&op)
	}

	return result
}
