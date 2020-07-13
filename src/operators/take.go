package operators

import "github.com/amao/toy-rxgo/src/base"

type takeSubscriber struct {
	*base.Subscriber
	valueCount uint
	count      uint
}

func newTakeSubscriber(destination base.SubscriberLike, count uint) takeSubscriber {
	newInstance := new(takeSubscriber)
	super := base.NewSubscriber(destination)
	newInstance.Subscriber = &super
	newInstance.count = count
	return *newInstance
}

func (t *takeSubscriber) Next(value interface{}) {
	total := t.count
	t.valueCount = t.valueCount + 1
	count := t.valueCount
	if count <= total {
		t.Destination.Next(value)
		if count == total {
			t.Destination.Complete()
			t.Unsubscribe()
		}
	}
}

type takeOperator struct {
	count uint
}

func newTakeOperator(count uint) takeOperator {
	newInstance := new(takeOperator)
	newInstance.count = count
	return *newInstance
}

func (t *takeOperator) Call(subscriber base.SubscriberLike, source base.Subscribable) base.SubscriptionLike {
	nts := newTakeSubscriber(subscriber, t.count)
	return source.Subscribe(&nts)
}

func Take(count uint) base.OperatorFunction {
	result := func(source base.Observable) base.Observable {
		op := newTakeOperator(count)
		return source.Lift(&op)
	}

	return result
}
