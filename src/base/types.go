package base

type OperatorFunction = func(Observable) Observable
type UnaryFunction = func(interface{}) interface{}

type Unsubscribable interface {
	Unsubscribe()
}

type Observer interface {
	Next(value interface{})
	Error(error)
	Complete()
}

type Operator interface {
	Call(subscriber *Subscriber, source Observable) Unsubscribable
}

type Subscribable interface {
	Subscribe(func(interface{}), func(error), func()) Unsubscribable
}

type TimestampProviderAndSchedulerLike interface {
	Now() float64
	Schedule(work func(scheduler SchedulerAction, state interface{}), delay float64, state interface{}) Unsubscribable
}

type SchedulerAction interface {
	Schedule(state interface{}, delay float64) Unsubscribable
	Unsubscribable
}
