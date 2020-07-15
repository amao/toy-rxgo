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

func (o *OuterSubscriber) NotifyNext(outerValue interface{}, innerValue interface{}, outerIndex interface{}, innerIndex interface{}, innerSub InnerSubscriber) {
	o.Destination.Next(innerValue)
}

func (o *OuterSubscriber) NotifyError(err error) {
	o.Destination.Error(err)
}

func (o *OuterSubscriber) NotifyComplete(innerSub InnerSubscriber) {
	o.Destination.Complete()
}
