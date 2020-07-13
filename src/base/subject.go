package base

type Subject struct {
	*Observable
	observers   []Observer
	closed      bool
	isStopped   bool
	hasError    bool
	thrownError interface{}
}

func NewSubject() *Subject {
	newInstance := new(Subject)
	self := NewObservable(nil)
	newInstance.Observable = &self
	newInstance.Observable.subscribe = func(subscriber SubscriberLike) SubscriptionLike {
		if newInstance.closed {
			panic("Object Unsubscribed Error")
		} else if newInstance.hasError {
			subscriber.Error(newInstance.thrownError.(error))
			subscription := NewSubscription(nil)
			return &subscription
		} else if newInstance.isStopped {
			subscriber.Complete()
			subscription := NewSubscription(nil)
			return &subscription
		} else {
			newInstance.observers = append(newInstance.observers, subscriber)
			return NewSubjectSubscription(*newInstance, subscriber)
		}
	}
	return newInstance
}

func (s *Subject) Next(value interface{}) {
	if s.closed {
		panic("Object Unsubscribed Error")
	}

	if !s.isStopped {
		copy := s.observers
		for _, observer := range copy {
			observer.Next(value)
		}
	}
}

func (s *Subject) Error(err error) {
	if s.closed {
		panic("Object UnSubscribed Error")
	}

	s.hasError = true
	s.thrownError = err
	s.isStopped = true
	copy := s.observers
	for _, observer := range copy {
		observer.Error(err)
	}
	s.observers = []Observer{}
}

func (s *Subject) Complete() {
	if s.closed {
		panic("Object Unsubscribed Error")
	}

	s.isStopped = true
	copy := s.observers
	for _, observer := range copy {
		observer.Complete()
	}

	s.observers = []Observer{}
}

func (s *Subject) Unsubscribe() {
	s.isStopped = true
	s.closed = true
	s.observers = nil
}

func (s *Subject) Lift(operator Operator) Subscribable {
	anonymouseSubject := NewAnonymousSubject(s, s)
	anonymouseSubject.operator = operator
	return anonymouseSubject
}

func (s *Subject) AsObservable() Observable {
	observable := NewObservable(nil)
	observable.source = s
	return observable
}

type AnonymousSubject struct {
	*Subject
	source      Subscribable
	destination Observer
}

func NewAnonymousSubject(destination Observer, source Subscribable) AnonymousSubject {
	newInstance := new(AnonymousSubject)
	subject := NewSubject()
	newInstance.Subject = subject

	newInstance.destination = destination
	newInstance.source = source
	return *newInstance
}

func (a *AnonymousSubject) Next(value interface{}) {
	if a.destination != nil && a.destination.Next != nil {
		a.destination.Next(value)
	}
}

func (a *AnonymousSubject) Error(err error) {
	if a.destination != nil && a.destination.Error != nil {
		a.destination.Error(err)
	}
}

func (a *AnonymousSubject) Complete() {
	if a.destination != nil && a.destination.Complete != nil {
		a.destination.Complete()
	}
}
