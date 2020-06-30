package operators

import (
	"github.com/amao/toy-rxgo/src/base"
	"github.com/amao/toy-rxgo/src/util"
)

type switchMapSubscriber struct {
	outerSubscriber   base.OuterSubscriber
	innerSubscription *base.Subscriber
	subscriber        base.Subscriber
	project           func(value interface{}) base.Subscribable
}

func newSwitchMapSubscriber(destination base.Subscriber, project func(interface{}) base.Subscribable) switchMapSubscriber {
	newInstance := switchMapSubscriber{}
	newInstance.subscriber = base.NewSubscriber(destination.Next, destination.Error, destination.Complete)
	newInstance.outerSubscriber = base.NewOuterSubscriber(destination)
	newInstance.project = project
	return newInstance
}

func (s switchMapSubscriber) innerSub(result base.Subscribable, value interface{}) {
	innerSubscription := s.innerSubscription
	if innerSubscription != nil {
		innerSubscription.Unsubscribe()
	}

	innerSubscriber := base.NewInnerSubscriber(s.outerSubscriber, value)
	destination := s.subscriber.Subscription
	destination.Add(innerSubscriber)
	in := util.SubscribeToResult(s.outerSubscriber, result, value, innerSubscriber).(base.Subscriber)
	s.innerSubscription = &in
}

func (s switchMapSubscriber) Next(value interface{}) {
	result := s.project(value)

	s.innerSub(result, value)
}

func (s switchMapSubscriber) Error(err error) {
	s.subscriber.Error(err)
}

func (s switchMapSubscriber) Complete() {
	innerSubscription := s.innerSubscription
	if innerSubscription == nil || innerSubscription.Subscription.Closed {
		s.subscriber.Complete()
	}
	s.Unsubscribe()
}

func (s switchMapSubscriber) Unsubscribe() {
	s.innerSubscription = nil
}

func (s switchMapSubscriber) notifyComplete(innerSub base.Subscription) {
	destination := s.subscriber.Subscription
	destination.Remove(innerSub)
	s.innerSubscription = nil
	if s.subscriber.IsStopped {
		s.subscriber.Complete()
	}
}

func (s switchMapSubscriber) notifyNext(innerValue interface{}) {
	s.subscriber.Destination.Next(innerValue)
}

type switchMapOperator struct {
	Project func(interface{}) base.Subscribable
}

func newSwitchMapOperator(project func(interface{}) base.Subscribable) switchMapOperator {
	newInstance := switchMapOperator{}
	newInstance.Project = project
	return newInstance
}

func (m switchMapOperator) Call(subscriber base.Subscriber, source base.Observable) base.Unsubscribable {
	nsms := newSwitchMapSubscriber(subscriber, m.Project)
	return source.Subscribe(nsms.Next, nsms.Error, nsms.Complete)
}

func SwitchMap(project func(value interface{}) base.Subscribable) base.OperatorFunction {
	result := func(source base.Observable) base.Observable {
		return source.Lift(newSwitchMapOperator(project))
	}

	return result
}
