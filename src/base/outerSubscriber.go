package base

type OuterSubscriber struct {
	*Subscriber
}

func NewOuterSubscriber() OuterSubscriber {
	newInstance := new(OuterSubscriber)
	parentSubscriber := NewSubscriber()
	newInstance.Subscriber = &parentSubscriber
	return *newInstance
}

func (o *OuterSubscriber) NotifyNext(innerValue interface{}) {
	o.Destination.Next(innerValue)
}

func (o *OuterSubscriber) NotifyError(err error) {
	o.Destination.Error(err)
}

func (o *OuterSubscriber) NotifyComplete(innerSub InnerSubscriber) {
	o.Destination.Complete()
}
