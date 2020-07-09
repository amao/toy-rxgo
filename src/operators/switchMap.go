package operators

import (
	"github.com/amao/toy-rxgo/src/base"
)

type switchMapSubscriber struct {
	parent            base.Subscriber
	fn                func(interface{}) base.Subscribable
	innerSubscription base.SubscriptionLike
	index             int
}

func newSwitchMapSubscriber(subscriber base.SubscriberLike, fn func(interface{}) base.Subscribable) switchMapSubscriber {
	newInstance := new(switchMapSubscriber)
	newInstance.index = 1
	newInstance.parent = base.NewSubscriber(subscriber)
	newInstance.fn = fn
	sp := base.NewSubscription(nil)
	newInstance.innerSubscription = &sp
	subscriber.Add(newInstance.innerSubscription)
	return *newInstance
}

func (s *switchMapSubscriber) Next(value interface{}) {
	innerObservable := s.fn(value)
	s.index++

	if s.innerSubscription != nil {
		s.innerSubscription.Unsubscribe()
	}

	s.innerSubscription = innerObservable.Subscribe(
		func(value interface{}) {
			s.parent.Next(value)
		},
		func(err error) {
			s.parent.Error(err)
		},
		func() {
			s.index--
			if s.index <= 0 {
				s.parent.Complete()
			}
		},
	)
}

func (s *switchMapSubscriber) Error(err error) {
	s.parent.Error(err)
}

func (s *switchMapSubscriber) Complete() {
	s.index = s.index - 1
	if s.index <= 0 {
		s.parent.Complete()
	}
}

type switchMapOperator struct {
	project func(interface{}) base.Subscribable
}

func newSwitchMapOperator(project func(interface{}) base.Subscribable) switchMapOperator {
	newInstance := new(switchMapOperator)
	newInstance.project = project
	return *newInstance
}

func (s *switchMapOperator) Call(subscriber base.SubscriberLike, source base.Observable) base.SubscriptionLike {
	nsms := newSwitchMapSubscriber(subscriber, s.project)
	return source.Subscribe(nsms.Next, nsms.Error, nsms.Complete)
}

func SwitchMap(fn func(interface{}) base.Subscribable) base.OperatorFunction {
	return func(source base.Observable) base.Observable {
		op := newSwitchMapOperator(fn)
		ob := source.Lift(&op)
		return ob
	}
}

func SwitchMap_(transformFn func(interface{}) base.Subscribable) base.OperatorFunction {
	result := func(inObservable base.Observable) base.Observable {
		outObservable := base.NewObservable(func(outObserver base.SubscriberLike) base.SubscriptionLike {
			inner_subscription := base.NewSubscription(func() {})
			var innerSubscription base.SubscriptionLike = &inner_subscription
			active := 1
			inObserver := base.NewSubscriber(
				func(x interface{}) {
					inner := transformFn(x)
					active++
					innerObserver := base.NewSubscriber(
						func(y interface{}) {
							outObserver.Next(y)
						},
						func(err error) {
							outObserver.Error(err)
						},
						func() {
							active = active - 1
							if active <= 0 {
								outObserver.Complete()
							}
						},
					)
					innerSubscription.Unsubscribe()
					is := inner.Subscribe(innerObserver).(*base.Subscriber)
					innerSubscription = is
				},
				func(err error) {
					outObserver.Error(err)
				},
				func() {
					active = active - 1
					if active <= 0 {
						outObserver.Complete()
					}
				},
			)
			subscription := inObserver.Subscription
			subscription.Add(innerSubscription)
			s := inObservable.Subscribe(inObserver)
			subscription.Add(s)
			return subscription
		})
		return outObservable
	}

	return result
}
