package operators

import (
	"github.com/amao/toy-rxgo/src/base"
)

type switchMapSubscriber struct {
	parent            base.Subscriber
	fn                func(interface{}) base.Subscribable
	innerSubscription base.Unsubscribable
}

func newSwitchMapSubscriber(subscriber *base.Subscriber, fn func(interface{}) base.Subscribable) switchMapSubscriber {
	newInstance := new(switchMapSubscriber)
	newInstance.parent = base.NewSubscriber(subscriber.Next, subscriber.Error, subscriber.Complete)
	newInstance.fn = fn
	sp := base.NewSubscription()
	newInstance.innerSubscription = &sp
	subscriber.Subscription.Add(newInstance.innerSubscription)
	return *newInstance
}

func (s *switchMapSubscriber) Next(value interface{}) {
	innerObservable := s.fn(value)

	if s.innerSubscription != nil {
		s.innerSubscription.Unsubscribe()
	}

	s.innerSubscription = innerObservable.Subscribe(
		func(value interface{}) {
			s.parent.Destination.Next(value)
		},
		nil,
		nil,
	)
}

func (s switchMapSubscriber) Error(err error) {
	s.parent.Error(err)
}

func (s switchMapSubscriber) Complete() {
	s.parent.Complete()
}

type switchMapOperator struct {
	project func(interface{}) base.Subscribable
}

func newSwitchMapOperator(project func(interface{}) base.Subscribable) switchMapOperator {
	newInstance := new(switchMapOperator)
	newInstance.project = project
	return *newInstance
}

func (m *switchMapOperator) Call(subscriber *base.Subscriber, source base.Observable) base.Unsubscribable {
	nsms := newSwitchMapSubscriber(subscriber, m.project)
	return source.Subscribe(nsms.Next, nsms.Error, nsms.Complete)
}

func SwitchMap(fn func(interface{}) base.Subscribable) base.OperatorFunction {
	return func(source base.Observable) base.Observable {
		op := newSwitchMapOperator(fn)
		return source.Lift(&op)
	}
}
