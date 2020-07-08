package base

type Subscriber struct {
	Destination Observer
	*Subscription
	IsStopped bool
}

func NewSubscriber(args ...interface{}) Subscriber {
	newInstance := new(Subscriber)
	newInstance.Subscription = new(Subscription)
	switch len(args) {
	case 0:
		emptyO := NewEmptyObserver()
		newInstance.Destination = &emptyO
	case 1:
		if obv, ok := args[0].(Observer); ok {
			newInstance.Destination = obv
		} else if nextFn, ok := args[0].(func(interface{})); ok {
			safeSub := newSafeSubscriber(*newInstance, nextFn, nil, nil)
			newInstance.Destination = &safeSub
		} else {
			emptyO := NewEmptyObserver()
			newInstance.Destination = &emptyO
		}
	case 2:
		if nextFn, ok := args[0].(func(interface{})); ok {
			if errFn, ok := args[1].(func(error)); ok {
				safeSub := newSafeSubscriber(*newInstance, nextFn, errFn, nil)
				newInstance.Destination = &safeSub
				break
			}
		}
		emptyO := NewEmptyObserver()
		newInstance.Destination = &emptyO
	case 3:
		if nextFn, ok := args[0].(func(interface{})); ok {
			if errFn, ok := args[1].(func(error)); ok {
				if completeFn, ok := args[2].(func()); ok {
					safeSub := newSafeSubscriber(*newInstance, nextFn, errFn, completeFn)
					newInstance.Destination = &safeSub
					break
				}
			}
		}
		emptyO := NewEmptyObserver()
		newInstance.Destination = &emptyO
	default:
		safeSub := newSafeSubscriber(*newInstance, args[0].(func(interface{})), args[1].(func(error)), args[2].(func()))
		newInstance.Destination = &safeSub
	}

	return *newInstance
}

func (s *Subscriber) _Next(value interface{}) {
	s.Destination.Next(value)
}

func (s *Subscriber) _Error(err error) {
	s.Destination.Error(err)
	s.Unsubscribe()
}

func (s *Subscriber) _Complete() {
	s.Destination.Complete()
	s.Unsubscribe()
}

func (s *Subscriber) Next(value interface{}) {
	if !s.IsStopped {
		s._Next(value)
	}
}

func (s *Subscriber) Error(err error) {
	if !s.IsStopped {
		s.IsStopped = true
		s._Error(err)
	}
}

func (s *Subscriber) Complete() {
	if !s.IsStopped {
		s.IsStopped = true
		s._Complete()
	}
}

func (s *Subscriber) Closed() bool {
	return s.closed
}

func (s *Subscriber) Unsubscribe() {
	if s.closed {
		return
	}

	s.IsStopped = true
	s.Subscription.Unsubscribe()
}

type safeSubscriber struct {
	subscriber *Subscriber
	next       func(value interface{})
	err        func(e error)
	complete   func()
}

func newSafeSubscriber(parent Subscriber, next func(value interface{}), err func(e error), complete func()) safeSubscriber {
	newInstance := new(safeSubscriber)
	subscriber := parent
	newInstance.subscriber = &subscriber
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

	return *newInstance

}

func (s *safeSubscriber) Next(value interface{}) {
	if !s.subscriber.IsStopped {
		s.next(value)
	}
}

func (s *safeSubscriber) Error(err error) {
	if !s.subscriber.IsStopped {
		s.err(err)
	}
}

func (s *safeSubscriber) Complete() {
	if !s.subscriber.IsStopped {
		s.complete()
	}
}

func (s *safeSubscriber) Closed() bool {
	return s.Closed()
}

func (s *safeSubscriber) Unsubscribe() {
	s.subscriber.Unsubscribe()
}
