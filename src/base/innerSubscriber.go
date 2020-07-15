package base

type InnerSubscriber struct {
	*Subscriber
	parent     OuterSubscriberLike
	outerValue interface{}
	outerIndex interface{}
}

func NewInnerSubscriber(parent OuterSubscriberLike, outerValue interface{}, outerIndex interface{}) InnerSubscriber {
	newInstance := new(InnerSubscriber)
	parentSubscriber := NewSubscriber()
	newInstance.Subscriber = &parentSubscriber
	newInstance.parent = parent
	newInstance.outerIndex = outerIndex
	newInstance.outerValue = outerValue

	newInstance.SetInnerNext(func(value interface{}) {
		newInstance.parent.NotifyNext(newInstance.outerValue, value, newInstance.outerIndex, -1, *newInstance)
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
