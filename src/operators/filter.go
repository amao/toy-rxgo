package operators

import "github.com/amao/toy-rxgo/src/base"

// another style of writting a operator
func filter(conditionFn func(interface{}) bool) func(base.Observable) base.Observable {
	result := func(inObservable base.Observable) base.Observable {
		outObservable := base.NewObservable(func(outObserver base.Observer) base.Unsubscribable {
			inObserver := base.NewSubscriber(func(x interface{}) {
				if conditionFn(x) {
					outObserver.Next(x)
				}
			}, func(e error) {
				outObserver.Error(e)
			}, func() {})

			return inObservable.Subscribe(inObserver.Next, inObserver.Error, inObserver.Complete)
		})
		return outObservable
	}

	return result
}

type filterSubscriber struct {
	subscriber base.Subscriber
	predicate  func(interface{}) bool
}

func newFilterSubscriber(destination base.Subscriber, predicate func(interface{}) bool) filterSubscriber {
	newInstance := filterSubscriber{}
	newInstance.subscriber = base.NewSubscriber(destination.Next, destination.Error, destination.Complete)
	newInstance.predicate = predicate
	return newInstance
}

func (s filterSubscriber) Next(value interface{}) {
	if s.predicate(value) {
		s.subscriber.Next(value)
	}
}

func (s filterSubscriber) Error(err error) {
	s.subscriber.Error(err)
}

func (s filterSubscriber) Complete() {
	s.subscriber.Complete()
}

type filterOperator struct {
	predicate func(interface{}) bool
}

func newFilterOperator(predicate func(interface{}) bool) filterOperator {
	newInstance := filterOperator{}
	newInstance.predicate = predicate
	return newInstance
}

func (m filterOperator) Call(subscriber base.Subscriber, source base.Observable) base.Unsubscribable {
	nfs := newFilterSubscriber(subscriber, m.predicate)
	return source.Subscribe(nfs.Next, nfs.Error, nfs.Complete)
}

func Filter(predicate func(interface{}) bool) base.OperatorFunction {
	result := func(source base.Observable) base.Observable {
		return source.Lift(newFilterOperator(predicate))
	}

	return result
}
