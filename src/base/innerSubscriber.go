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
	return *newInstance
}

func (i *InnerSubscriber) Next(value interface{}) {
	i.parent.NotifyNext(value)
}

func (i *InnerSubscriber) Error(err error) {
	i.parent.NotifyError(err)
	i.Unsubscribe()
}

func (i *InnerSubscriber) Complete() {
	i.parent.NotifComplete(*i)
	i.Unsubscribe()
}
