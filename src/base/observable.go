package base

type Observable struct {
	subscribe func(SubscriberLike) SubscriptionLike
	source    Subscribable
	operator  Operator
}

func NewObservable(subscribe interface{}) Observable {
	observable := new(Observable)
	if subscribe != nil {
		observable.subscribe = subscribe.(func(SubscriberLike) SubscriptionLike)
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
	observable.source = o
	observable.operator = operator
	return *observable
}

func (o *Observable) Subscribe(args ...interface{}) SubscriptionLike {
	sink := toSubscriber(args...)

	operator := o.operator

	var subscription SubscriptionLike

	if operator != nil {
		subscription = operator.Call(sink, o.source)
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
		if realType, ok := args[0].(func(interface{})); ok {
			result := NewSubscriber(realType)
			return &result
		}
		result := args[0].(SubscriberLike)
		return result
	default:
		result := NewSubscriber(args[0].(func(interface{})), args[1].(func(error)), args[2].(func()))
		return &result
	}
}
