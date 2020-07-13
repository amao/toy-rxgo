package base

type BehaviorSubject struct {
	*Subject
	super *Subject
	value interface{}
}

func NewBehaviorSubject(value interface{}) *BehaviorSubject {
	newInstance := new(BehaviorSubject)
	self := NewSubject()
	super := NewSubject()
	newInstance.Subject = self
	newInstance.super = super
	newInstance.value = value
	newInstance.subscribe = func(subscriber SubscriberLike) SubscriptionLike {
		subscription := super.subscribe(subscriber)
		if subscription != nil && !subscription.Closed() {
			subscriber.Next(newInstance.value)
		}
		return subscription
	}
	return newInstance
}

func (b *BehaviorSubject) Value() interface{} {
	if b.hasError {
		panic(b.thrownError)
	} else if b.closed {
		panic("Object Unsubscribed Error")
	} else {
		return b.value
	}
}

func (b *BehaviorSubject) Next(value interface{}) {
	b.value = value
	b.super.Next(b.value)
}
