package base

type OperatorFunction = func(Observable) Observable
type UnaryFunction = func(interface{}) interface{}

type Unsubscribable interface {
	Unsubscribe()
}

type SubscriptionLike interface {
	Unsubscribable
	Closed() bool
}

type Observer interface {
	Next(value interface{})
	Error(error)
	Complete()
}

type Operator interface {
	Call(subscriber *Subscriber, source Observable) SubscriptionLike
}

type Subscribable interface {
	Subscribe(args ...interface{}) SubscriptionLike
}
