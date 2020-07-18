package operators

import (
	"github.com/amao/toy-rxgo/src/base"
	"github.com/amao/toy-rxgo/src/util"
)

type auditSubscriber struct {
	*base.OuterSubscriber
	value            interface{}
	hasValue         bool
	throttled        base.SubscriptionLike
	durationSelector func(interface{}) base.Subscribable
}

func newAuditSubscriber(destination base.SubscriberLike, durationSelector func(interface{}) base.Subscribable) auditSubscriber {
	newInstance := new(auditSubscriber)

	subscriber := base.NewSubscriber(destination)
	self := base.NewOuterSubscriber()
	self.Subscriber = &subscriber
	newInstance.OuterSubscriber = &self
	newInstance.Destination.Add(newInstance)

	newInstance.durationSelector = durationSelector

	newInstance.SetInnerNext(func(value interface{}) {
		newInstance.value = value
		newInstance.hasValue = true
		if newInstance.throttled == nil {
			var duration base.Subscribable
			durationSelector := newInstance.durationSelector
			duration = durationSelector(value)

			innerSubscriber := base.NewInnerSubscriber(newInstance, nil, nil)
			innerSubscription := duration.Subscribe(innerSubscriber)
			if innerSubscription == nil || innerSubscription.Closed() {
				newInstance.clearThrottle()
			} else {
				newInstance.throttled = innerSubscription
				newInstance.Add(newInstance.throttled)
			}
		}

	})

	return *newInstance
}

func (a *auditSubscriber) clearThrottle() {
	value := a.value
	hasValue := a.hasValue
	throttled := a.throttled

	if throttled != nil {
		a.Remove(throttled)
		a.throttled = nil
		throttled.Unsubscribe()
	}

	if hasValue {
		a.value = nil
		a.hasValue = false
		a.Destination.Next(value)
	}
}

func (a *auditSubscriber) NotifyNext(outerValue interface{}, innerValue interface{}, outerIndex interface{}, innerIndex interface{}, innerSub base.InnerSubscriber) {

	a.clearThrottle()
}

func (a *auditSubscriber) NotifyComplete(innerSub base.InnerSubscriber) {
	a.clearThrottle()
}

type auditOperator struct {
	durationSelector func(value interface{}) base.Subscribable
}

func newAuditOperator(durationSelector func(value interface{}) base.Subscribable) auditOperator {
	newInstance := new(auditOperator)
	newInstance.durationSelector = durationSelector
	return *newInstance
}

func (a *auditOperator) Call(subscriber base.SubscriberLike, source base.Subscribable) base.SubscriptionLike {
	nas := newAuditSubscriber(subscriber, a.durationSelector)
	return source.Subscribe(nas)
}

func Audit(durationSelector func(value interface{}) base.Subscribable) base.OperatorFunction {
	result := func(source base.Observable) base.Observable {
		op := newAuditOperator(durationSelector)
		return util.Lift(source, &op)
	}
	return result
}
