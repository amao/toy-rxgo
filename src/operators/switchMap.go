package operators

import (
	"github.com/amao/toy-rxgo/src/base"
)

type switchMapSubscriber struct {
	*base.OuterSubscriber
	super             *base.OuterSubscriber
	project           func(interface{}) base.Subscribable
	innerSubscription base.SubscriptionLike
}

func newSwitchMapSubscriber(destination base.SubscriberLike, project func(interface{}) base.Subscribable) switchMapSubscriber {
	newInstance := new(switchMapSubscriber)

	selfSubscriber := base.NewSubscriber(destination)
	selfOuterSubscriber := new(base.OuterSubscriber)
	newInstance.OuterSubscriber = selfOuterSubscriber
	newInstance.Subscriber = &selfSubscriber
	newInstance.Destination.Add(newInstance)

	superSubscriber := base.NewSubscriber(destination)
	superOuterSubscriber := new(base.OuterSubscriber)
	newInstance.super = superOuterSubscriber
	newInstance.super.Subscriber = &superSubscriber

	newInstance.project = project
	//sp := base.NewSubscription(nil)
	//newInstance.innerSubscription = &sp
	//destination.Add(newInstance.innerSubscription)

	newInstance.SetInnerNext(func(value interface{}) {
		result := newInstance.project(value)
		newInstance._innerSub(result, value)
	})
	newInstance.SetInnerComplete(func() {
		if newInstance.innerSubscription == nil || newInstance.innerSubscription.Closed() {
			newInstance.super.CallInnerComplete()
		}
		newInstance.Unsubscribe()
	})
	newInstance.SetInnerUnsubscribe(func() {
		newInstance.innerSubscription = nil
	})
	return *newInstance
}

func (s *switchMapSubscriber) _innerSub(result base.Subscribable, value interface{}) {

	if s.innerSubscription != nil {
		s.innerSubscription.Unsubscribe()
	}

	innerSubscriber := base.NewInnerSubscriber(s)
	//destination := s.Destination
	//destination.Add(innerSubscriber)

	if innerSubscriber.Closed() {
		return
	}

	s.innerSubscription = result.Subscribe(&innerSubscriber)
}

func (s *switchMapSubscriber) NotifyNext(innerValue interface{}) {
	s.Destination.Next(innerValue)
}

func (s *switchMapSubscriber) NotifComplete(innerSub base.InnerSubscriber) {
	destination := s.Destination
	destination.Remove(innerSub)
	s.innerSubscription = nil
	if s.IsStopped {
		s.super.CallInnerComplete()
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
	return source.Subscribe(nsms)
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
