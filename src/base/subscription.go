package base

import (
	"reflect"
)

type Subscription struct {
	closed          bool //default value = false
	parentOrParents interface{}
	subscriptions   []Unsubscribable
	unsubscribe     func()
}

func NewSubscription(unsubscribe interface{}) Subscription {
	newInstance := new(Subscription)
	if unsubscribe != nil {
		newInstance.unsubscribe = unsubscribe.(func())
		newInstance.closed = false
	}
	return *newInstance
}

func (s *Subscription) Closed() bool {
	return s.closed
}

func (s *Subscription) Add(teardown SubscriptionLike) Unsubscribable {

	if reflect.DeepEqual(teardown, s) || teardown.Closed() {
		return teardown
	} else if s.Closed() {
		teardown.Unsubscribe()
		return teardown
	}

	if s.subscriptions == nil {
		s.subscriptions = []Unsubscribable{teardown}
	} else {
		s.subscriptions = append(s.subscriptions, teardown)
	}

	return teardown
}

func (s *Subscription) Remove(subscription Subscription) {
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

	if s.unsubscribe != nil {
		s.unsubscribe()
	}

	for _, sub := range s.subscriptions {
		sub.Unsubscribe()
	}
}
