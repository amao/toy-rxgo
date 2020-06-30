package base

type Subscriber struct {
	Destination  Observer
	Subscription Subscription
	IsStopped    bool
}

func NewSubscriber(args ...interface{}) Subscriber {
	newInstance := Subscriber{}
	switch len(args) {
	case 0:
		newInstance.Destination = NewEmptyObserver()
	case 1:
		if obv, ok := args[0].(Observer); ok {
			newInstance.Destination = obv
		} else if nextFn, ok := args[0].(func(interface{})); ok {
			newInstance.Destination = newSafeSubscriber(nextFn, nil, nil)
		} else {
			newInstance.Destination = NewEmptyObserver()
		}
	case 2:
		if nextFn, ok := args[0].(func(interface{})); ok {
			if errFn, ok := args[1].(func(error)); ok {
				newInstance.Destination = newSafeSubscriber(nextFn, errFn, nil)
				break
			}
		}
		newInstance.Destination = NewEmptyObserver()
	case 3:
		if nextFn, ok := args[0].(func(interface{})); ok {
			if errFn, ok := args[1].(func(error)); ok {
				if completeFn, ok := args[2].(func()); ok {
					newInstance.Destination = newSafeSubscriber(nextFn, errFn, completeFn)
					break
				}
			}
		}
		newInstance.Destination = NewEmptyObserver()
	default:
		newInstance.Destination = newSafeSubscriber(args[0].(func(interface{})), args[1].(func(error)), args[2].(func()))
	}

	return newInstance
}

func (s Subscriber) _Next(value interface{}) {
	s.Destination.Next(value)
}

func (s Subscriber) _Error(err error) {
	s.Destination.Error(err)
	s.Unsubscribe()
}

func (s Subscriber) _Complete() {
	s.Destination.Complete()
	s.Unsubscribe()
}

func (s Subscriber) Next(value interface{}) {
	if !s.IsStopped {
		s._Next(value)
	}
}

func (s Subscriber) Error(err error) {
	if !s.IsStopped {
		s.IsStopped = true
		s._Error(err)
	}
}

func (s Subscriber) Complete() {
	if !s.IsStopped {
		s.IsStopped = true
		s._Complete()
	}
}

func (s Subscriber) Unsubscribe() {
	if s.Subscription.Closed {
		return
	}

	s.IsStopped = true
	s.Subscription.Unsubscribe()
}

type safeSubscriber struct {
	subscriber Subscriber
	next       func(value interface{})
	err        func(e error)
	complete   func()
}

func newSafeSubscriber(next func(value interface{}), err func(e error), complete func()) safeSubscriber {
	newInstance := safeSubscriber{}
	newInstance.subscriber = NewSubscriber()
	newInstance.next = next
	if err != nil {
		newInstance.err = err
	} else {
		newInstance.err = func(err error) {}
	}
	if complete != nil {
		newInstance.complete = complete
	} else {
		newInstance.complete = func() {}
	}

	return newInstance

}

func (s safeSubscriber) Next(value interface{}) {
	if !s.subscriber.IsStopped {
		s.next(value)
	}
}

func (s safeSubscriber) Error(err error) {
	if !s.subscriber.IsStopped {
		s.err(err)
	}
}

func (s safeSubscriber) Complete() {
	if !s.subscriber.IsStopped {
		s.complete()
	}
}

func (s safeSubscriber) Unsubscribe() {
	s.subscriber.Unsubscribe()
}
