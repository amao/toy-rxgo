package base

type Observable struct {
	subscribe func(Observer) Unsubscribable
	source    *Observable
	operator  Operator
}

func NewObservable(subscribe func(Observer) Unsubscribable) Observable {
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

func (o *Observable) Lift(operator Operator) Observable {
	observable := new(Observable)
	observable.subscribe = func(subscriber Observer) Unsubscribable {
		if o.source != nil {
			o.source.Subscribe(subscriber.Next, subscriber.Error, subscriber.Complete)
		}
		subscription := NewSubscription()
		return &subscription
	}
	observable.source = o
	observable.operator = operator
	return *observable
}

func (o *Observable) Subscribe(next func(interface{}), err func(err error), complete func()) Unsubscribable {
	sink := NewSubscriber(next, err, complete)

	operator := o.operator

	var subscription Unsubscribable

	if operator != nil {
		subscription = operator.Call(&sink, *o.source)
	} else {
		subscription = o.subscribe(&sink)
	}

	sink.Subscription.Add(subscription)

	return &sink
}
