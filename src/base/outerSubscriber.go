package base

type OuterSubscriber struct {
	subscriber Subscriber
}

func NewOuterSubscriber(destination Subscriber) OuterSubscriber {
	newInstance := OuterSubscriber{}
	subscriber := NewSubscriber(destination.Next, destination.Error, destination.Complete)

	newInstance.subscriber = subscriber
	return newInstance
}

func (o OuterSubscriber) notifyNext(outerValue interface{}, innerValue interface{}) {
	o.subscriber.Destination.Next(innerValue)
}

func (o OuterSubscriber) notifyError(err error) {
	o.subscriber.Destination.Error(err)
}

func (o OuterSubscriber) notifyComplete() {
	o.subscriber.Destination.Complete()
}

func (o OuterSubscriber) Next(innerValue interface{}) {
	o.subscriber.Next(innerValue)
}

func (o OuterSubscriber) Error(err error) {
	o.subscriber.Error(err)
}

func (o OuterSubscriber) Complete() {
	o.subscriber.Complete()
}
