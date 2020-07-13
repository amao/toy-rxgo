package base

import "reflect"

type SubjectSubscription struct {
	*Subscription
	closed     bool
	subject    *Subject
	subscriber Observer
}

func NewSubjectSubscription(subject Subject, subscriber Observer) SubscriptionLike {
	newInstance := new(SubjectSubscription)
	self := NewSubscription(nil)
	newInstance.Subscription = &self
	newInstance.subject = &subject
	newInstance.subscriber = subscriber
	return newInstance
}

func (s *SubjectSubscription) Unsubscribe() {
	if s.closed {
		return
	}

	s.closed = true

	subject := s.subject
	observers := subject.observers

	s.subject = nil

	if observers != nil || len(observers) == 0 || subject.isStopped || subject.closed {
		return
	}

	var subscriberIndex = -1

	for index, subscriber := range observers {
		if reflect.DeepEqual(subscriber, s.subscriber) {
			subscriberIndex = index
		}
	}

	if subscriberIndex != -1 {
		result := observers[:subscriberIndex]
		result = append(result, observers[subscriberIndex:]...)
		observers = result
	}
}
