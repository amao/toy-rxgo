package base

type Observable struct {
	subscribe func(Observer) SubscriptionLike
	source    *Observable
	operator  Operator
}

func NewObservable(subscribe func(Observer) SubscriptionLike) Observable {
	observable := new(Observable)
	if subscribe != nil {
		observable.subscribe = subscribe
	}
	return *observable
}

func (o *Observable) Pipe(operations ...OperatorFunction) *Observable {
	if len(operations) == 0 {
		return o
	}

	var result Observable = *o

	for _, fn := range operations {
		result = fn(result)
	}
	return &result
}

func (o *Observable) Lift(operator Operator) *Observable {
	observable := new(Observable)
	observable.subscribe = func(subscriber Observer) SubscriptionLike {
		if o.source != nil {
			o.source.Subscribe(subscriber.Next, subscriber.Error, subscriber.Complete)
		}
		subscription := NewSubscription(nil)
		return &subscription
	}
	observable.source = o
	observable.operator = operator
	return observable
}

func (o *Observable) Subscribe(args ...interface{}) SubscriptionLike {
	sink := toSubscriber(args...)

	operator := o.operator

	var subscription SubscriptionLike

	if operator != nil {
		subscription = operator.Call(sink, *o.source)
	} else {
		subscription = o.subscribe(sink)
	}

	sink.Add(subscription)

	return sink
}

func toSubscriber(args ...interface{}) SubscriberLike {
	switch len(args) {
	case 0:
		result := NewSubscriber()
		return &result
	case 1:
		result := args[0].(SubscriberLike)
		return result
	default:
		result := NewSubscriber(args[0].(func(interface{})), args[1].(func(error)), args[2].(func()))
		return &result
	}
}
