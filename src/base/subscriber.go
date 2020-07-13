package base

type Subscriber struct {
	Destination SubscriberLike
	super       *Subscription
	*Subscription
	IsStopped bool
	_next     func(interface{})
	_error    func(err error)
	_complete func()
}

func NewSubscriber(args ...interface{}) Subscriber {
	newInstance := new(Subscriber)
	newInstance.Subscription = new(Subscription)
	newInstance.super = new(Subscription)
	switch len(args) {
	case 0:
		emptyObserver := NewEmptyObserver()
		newInstance.Destination = emptyObserver
	case 1:
		if obv, ok := args[0].(SubscriberLike); ok {
			newInstance.Destination = obv
			// New sub struct needs to be added Subscription of Destination
		} else if nextFn, ok := args[0].(func(interface{})); ok {
			safeSub := newSafeSubscriber(*newInstance, nextFn, nil, nil)
			newInstance.Destination = safeSub
		} else {
			emptyObserver := NewEmptyObserver()
			newInstance.Destination = emptyObserver
		}
	case 2:
		if nextFn, ok := args[0].(func(interface{})); ok {
			if errFn, ok := args[1].(func(error)); ok {
				safeSub := newSafeSubscriber(*newInstance, nextFn, errFn, nil)
				newInstance.Destination = safeSub
				break
			}
		}
		emptyO := NewEmptyObserver()
		newInstance.Destination = emptyO
	case 3:
		if nextFn, ok := args[0].(func(interface{})); ok {
			if errFn, ok := args[1].(func(error)); ok {
				if completeFn, ok := args[2].(func()); ok {
					safeSub := newSafeSubscriber(*newInstance, nextFn, errFn, completeFn)
					newInstance.Destination = safeSub
					break
				}
			}
		}
		emptyO := NewEmptyObserver()
		newInstance.Destination = emptyO
	default:
		safeSub := newSafeSubscriber(*newInstance, args[0].(func(interface{})), args[1].(func(error)), args[2].(func()))
		newInstance.Destination = safeSub
	}

	newInstance._next = func(value interface{}) {
		newInstance.Destination.Next(value)
	}

	newInstance._error = func(err error) {
		newInstance.Destination.Error(err)
		newInstance.Unsubscribe()
	}

	newInstance._complete = func() {
		newInstance.Destination.Complete()
		newInstance.Unsubscribe()
	}

	return *newInstance
}

func (s *Subscriber) SetInnerNext(_next func(interface{})) {
	s._next = _next
}

func (s *Subscriber) SetInnerError(_error func(error)) {
	s._error = _error
}

func (s *Subscriber) SetInnerComplete(_complete func()) {
	s._complete = _complete
}

func (s *Subscriber) CallInnerNext(value interface{}) {
	s._next(value)
}

func (s *Subscriber) CallInnerComplete() {
	s._complete()
}

func (s *Subscriber) Next(value interface{}) {
	if !s.IsStopped {
		s._next(value)
	}
}

func (s *Subscriber) Error(err error) {
	if !s.IsStopped {
		s.IsStopped = true
		s._error(err)
	}
}

func (s *Subscriber) Complete() {
	if !s.IsStopped {
		s.IsStopped = true
		s._complete()
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
	s.super.Unsubscribe()
}

type safeSubscriber struct {
	*Subscriber
	super     *Subscriber
	parent    *Subscriber
	IsStopped bool
	closed    bool
}

func newSafeSubscriber(parent Subscriber, next func(value interface{}), err func(e error), complete func()) SubscriberLike {
	newInstance := new(safeSubscriber)
	self := NewSubscriber()
	super := NewSubscriber()
	newInstance.Subscriber = &self
	newInstance.super = &super
	newInstance.parent = &parent
	newInstance._next = next
	if err != nil {
		newInstance._error = err
	} else {
		newInstance._error = func(err error) {}
	}
	if complete != nil {
		newInstance._complete = complete
	} else {
		newInstance._complete = func() {}
	}

	newInstance.SetInnerUnsubscribe(func() {
		if newInstance.parent != nil {
			newInstance.parent.Unsubscribe()
		}
	})

	return newInstance

}

func (s *safeSubscriber) Next(value interface{}) {
	if !s.IsStopped {
		s._next(value)
	}
}

func (s *safeSubscriber) Error(err error) {
	if !s.IsStopped {
		s._error(err)
		s.Unsubscribe()
	}
}

func (s *safeSubscriber) Complete() {
	if !s.IsStopped {
		s._complete()
		s.Unsubscribe()
	}
}

func (s *safeSubscriber) Closed() bool {
	return s.closed
}
