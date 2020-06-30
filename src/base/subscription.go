package base

import (
	"reflect"
)

type Subscription struct {
	Closed        bool
	subscriptions []Unsubscribable
	unsubscribe   func()
}

func NewSubscription(args ...func()) Subscription {
	newInstance := Subscription{}
	if len(args) != 0 {
		newInstance.unsubscribe = args[0]
		newInstance.Closed = false
	}
	return newInstance
}

func (s Subscription) Add(teardown Unsubscribable) Unsubscribable {
	var subscriptionLike Unsubscribable
	if subscription, ok := teardown.(Subscription); ok {
		subscriptionLike = subscription
		if reflect.DeepEqual(subscription, s) || subscription.Closed {
			return subscription
		} else if s.Closed {
			subscription.Unsubscribe()
			return subscription
		}
	}

	if subscriber, ok := teardown.(Subscriber); ok {
		subscriptionLike = subscriber
		if reflect.DeepEqual(subscriber, s) || subscriber.Subscription.Closed {
			return subscriber
		} else if s.Closed {
			subscriber.Unsubscribe()
			return subscriber
		}
	}

	subscriptions := s.subscriptions
	if subscriptions == nil {
		s.subscriptions = []Unsubscribable{subscriptionLike}
	} else {
		subscriptions = append(subscriptions, subscriptionLike)
	}

	return subscriptionLike
}

func (s Subscription) Remove(subscription Subscription) {
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

func (s Subscription) Unsubscribe() {
	if s.Closed {
		return
	}

	s.Closed = true

	if s.unsubscribe != nil {
		s.unsubscribe()
	}

	for _, sub := range s.subscriptions {
		sub.Unsubscribe()
	}
}
