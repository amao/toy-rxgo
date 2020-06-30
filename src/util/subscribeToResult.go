package util

import (
	"github.com/amao/toy-rxgo/src/base"
)

func SubscribeToResult(outerSubscriber base.OuterSubscriber, result base.Subscribable, outerValue interface{}, innerSubscriber base.InnerSubscriber) base.Unsubscribable {
	if innerSubscriber.Subscriber.Subscription.Closed {
		return nil
	}
	re := result.Subscribe(innerSubscriber.Next, innerSubscriber.Error, innerSubscriber.Complete)
	return re
}
