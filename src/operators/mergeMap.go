package operators

import (
	"fmt"

	"github.com/amao/toy-rxgo/src/base"
)

type mergeMapSubscriber struct {
	*base.Subscriber
	project func(interface{}) base.Subscribable
	index   int
}

func newMergeMapSubscriber(subscriber *base.Subscriber, project func(interface{}) base.Subscribable) mergeMapSubscriber {
	newInstance := new(mergeMapSubscriber)
	newInstance.index = 1
	s := base.NewSubscriber(subscriber)
	newInstance.Subscriber = &s
	newInstance.project = project
	return *newInstance
}

func (m *mergeMapSubscriber) Next(value interface{}) {
	fmt.Println("from value: ", value)
	innerObservable := m.project(value)
	m.index++

	sink := innerObservable.Subscribe(
		func(value interface{}) {
			fmt.Println("to value: ", value)
			m.Destination.Next(value)
		},
		func(err error) {
			m.Destination.Error(err)
		},
		func() {
			m.index--
			if m.index <= 0 {
				m.Destination.Complete()
			}
		},
	)
	m.Add(sink)
}

func (m *mergeMapSubscriber) Error(err error) {
	m.Destination.Error(err)
}

func (m *mergeMapSubscriber) Complete() {
	m.index--
	if m.index <= 0 {
		m.Destination.Complete()
	}
}

type mergeMapOperator struct {
	project func(interface{}) base.Subscribable
}

func newMergeMapOperator(project func(interface{}) base.Subscribable) mergeMapOperator {
	newInstance := new(mergeMapOperator)
	newInstance.project = project
	return *newInstance
}

func (m *mergeMapOperator) Call(subscriber *base.Subscriber, source base.Observable) base.Unsubscribable {
	nmms := newMergeMapSubscriber(subscriber, m.project)
	return source.Subscribe(nmms.Next, nmms.Error, nmms.Complete)
}

func MergeMap_(fn func(interface{}) base.Subscribable) base.OperatorFunction {
	return func(source base.Observable) base.Observable {
		op := newMergeMapOperator(fn)
		ob := source.Lift(&op)
		return *ob
	}
}

func MergeMap(transformFn func(interface{}) base.Subscribable) base.OperatorFunction {
	result := func(inObservable base.Observable) base.Observable {
		outObservable := base.NewObservable(func(outObserver base.Observer) base.Unsubscribable {
			subscription := base.NewSubscription()
			active := 0
			inObserver := base.NewSubscriber(
				func(x interface{}) {
					inner := transformFn(x)
					innerObserver := base.NewSubscriber(
						func(y interface{}) {
							outObserver.Next(y)
						},
						func(err error) {
							outObserver.Error(err)
						},
						func() {
							active = active - 1
							if active == 0 {
								outObserver.Complete()
							}
						},
					)
					active++
					subscription.Add(inner.Subscribe(innerObserver))
				},
				func(err error) {
					outObserver.Error(err)
				},
				func() {
					active = active - 1
					if active == 0 {
						outObserver.Complete()
					}
				},
			)
			active++
			s := inObservable.Subscribe(inObserver)
			subscription.Add(s)
			return &subscription
		})
		return outObservable
	}

	return result
}
