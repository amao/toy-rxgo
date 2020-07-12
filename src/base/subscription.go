package base

import (
	"reflect"
)

type Subscription struct {
	closed          bool //default value = false
	parentOrParents interface{}
	subscriptions   []Unsubscribable
	_unsubscribe    func()
}

func NewSubscription(unsubscribe interface{}) Subscription {
	newInstance := new(Subscription)
	if unsubscribe != nil {
		newInstance._unsubscribe = unsubscribe.(func())
		newInstance.closed = false
	}
	return *newInstance
}

func (s *Subscription) SetInnerUnsubscribe(innerUnsubscribe func()) {
	s._unsubscribe = innerUnsubscribe
}

func (s *Subscription) Closed() bool {
	return s.closed
}

func (s *Subscription) Add(teardown SubscriptionLike) Unsubscribable {
	if subscription, ok := teardown.(*Subscription); ok {
		if reflect.DeepEqual(subscription, s) || subscription.closed {
			return teardown
		} else if s.closed {
			teardown.Unsubscribe()
			return teardown
		}
	}

	if subscriber, ok := teardown.(*Subscriber); ok {
		if reflect.DeepEqual(subscriber.Subscription, s) || subscriber.closed {
			return teardown
		} else if s.closed {
			teardown.Unsubscribe()
			return teardown
		}
	}

	if s.subscriptions == nil {
		s.subscriptions = []Unsubscribable{teardown}
	} else {
		s.subscriptions = append(s.subscriptions, teardown)
	}

	return teardown
}

func (s *Subscription) Remove(subscription SubscriptionLike) {
	subscriptions := s.subscriptions

	if subscriptions != nil {
		index := -1

		// Find index
		for i := 0; i < len(subscriptions); i++ {
			if reflect.DeepEqual(subscriptions[i], subscription) {
				index = i
				break
			}
		}

		if index != -1 {
			//Delete index
			subscriptions[index] = subscriptions[len(subscriptions)-1]
			subscriptions[len(subscriptions)-1] = nil
			subscriptions = subscriptions[:len(subscriptions)-1]

		}
	}
}

func (s *Subscription) Unsubscribe() {
	if s.closed {
		return
	}

	s.closed = true

	if s._unsubscribe != nil {
		s._unsubscribe()
	}

	for _, sub := range s.subscriptions {
		if sub != nil {
			sub.Unsubscribe()
		}

	}
}
