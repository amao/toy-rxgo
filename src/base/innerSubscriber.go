package base

type InnerSubscriber struct {
	*Subscriber
	parent OuterSubscriberLike
}

func NewInnerSubscriber(parent OuterSubscriberLike) InnerSubscriber {
	newInstance := new(InnerSubscriber)
	parentSubscriber := NewSubscriber()
	newInstance.Subscriber = &parentSubscriber
	newInstance.parent = parent
	newInstance.SetInnerNext(func(value interface{}) {
		newInstance.parent.NotifyNext(value)
	})
	newInstance.SetInnerError(func(err error) {
		newInstance.parent.NotifyError(err)
		newInstance.Unsubscribe()
	})
	newInstance.SetInnerComplete(func() {
		newInstance.parent.NotifyComplete(*newInstance)
		newInstance.Unsubscribe()
	})
	return *newInstance
}
