package base

type InnerSubscriber struct {
	Subscriber Subscriber
	parent     OuterSubscriber
	outValue   interface{}
}

func NewInnerSubscriber(parent OuterSubscriber, outValue interface{}) InnerSubscriber {
	newInstance := InnerSubscriber{}
	subscriber := NewSubscriber()
	newInstance.Subscriber = subscriber
	newInstance.parent = parent
	newInstance.outValue = outValue
	return newInstance
}
func (i InnerSubscriber) Next(value interface{}) {
	i.parent.notifyNext(i.outValue, value)
}

func (i InnerSubscriber) Error(err error) {
	i.parent.notifyError(err)
	i.Subscriber.Unsubscribe()
}

func (i InnerSubscriber) Complete() {
	i.parent.notifyComplete()
	i.Subscriber.Unsubscribe()
}

func (i InnerSubscriber) Unsubscribe() {
	i.Subscriber.Unsubscribe()
}
