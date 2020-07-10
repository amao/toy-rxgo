package base

import "time"

type OperatorFunction = func(Observable) Observable
type UnaryFunction = func(interface{}) interface{}

type Unsubscribable interface {
	Unsubscribe()
}

type SubscriptionLike interface {
	Unsubscribable
	Closed() bool
	Add(SubscriptionLike) Unsubscribable
	Remove(SubscriptionLike)
}

type SubscriberLike interface {
	Unsubscribable
	Closed() bool
	Add(SubscriptionLike) Unsubscribable
	Remove(SubscriptionLike)
	Observer
}

type Observer interface {
	Next(interface{})
	Error(error)
	Complete()
}

type Operator interface {
	Call(subscriber SubscriberLike, source Observable) SubscriptionLike
}

type Subscribable interface {
	Subscribe(args ...interface{}) SubscriptionLike
}

type OuterSubscriberLike interface {
	NotifyNext(interface{})
	NotifyError(error)
	NotifComplete(innerSub InnerSubscriber)
}

type SchedulerAction interface {
	SubscriptionLike
	Schedule(state interface{}, delay uint) SubscriptionLike
}

type TimestampProvider interface {
	Now() time.Time
}

type SchedulerLike interface {
	TimestampProvider
	Schedule(scheduler SchedulerLike, work func(SchedulerAction, interface{}), delay uint, state interface{}) SubscriptionLike
}
