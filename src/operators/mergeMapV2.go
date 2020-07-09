package operators

import (
	"github.com/amao/toy-rxgo/src/base"
)

type mergeMapSubscriberV2 struct {
	*base.OuterSubscriber
	hasCompleted bool //default value false
	buffer       []interface{}
	active       int //default value 0
	index        int //default value 0
	project      func(interface{}) base.Subscribable
	concurrent   int //default int max
}

func newMergeMapSubscriberV2(destination base.SubscriberLike, project func(interface{}) base.Subscribable, concurrent int) mergeMapSubscriberV2 {
	newInstance := new(mergeMapSubscriberV2)
	subscriber := base.NewSubscriber(destination)
	outerSubscriber := base.NewOuterSubscriber()
	outerSubscriber.Subscriber = &subscriber
	newInstance.OuterSubscriber = &outerSubscriber
	newInstance.project = project
	newInstance.concurrent = concurrent
	return *newInstance
}

func (m *mergeMapSubscriberV2) Next(value interface{}) {
	if m.active < m.concurrent {
		m._tryNext(value)
	} else {
		m.buffer = append(m.buffer, value)
	}
}

func (m *mergeMapSubscriberV2) Error(err error) {
	m.Destination.Error(err)

}

func (m *mergeMapSubscriberV2) Complete() {
	m.hasCompleted = true
	if m.active == 0 && len(m.buffer) == 0 {
		m.Destination.Complete()
	}
	m.Unsubscribe()
}

func (m *mergeMapSubscriberV2) _tryNext(value interface{}) {
	var result base.Subscribable
	result = m.project(value)
	m.active++
	m._innerSub(result, value)
}

func (m *mergeMapSubscriberV2) _innerSub(ish base.Subscribable, value interface{}) {
	innerSubscriber := base.NewInnerSubscriber(m)
	destination := m.Subscription
	destination.Add(innerSubscriber)

	if innerSubscriber.Closed() {
		return
	}

	ish.Subscribe(&innerSubscriber)
}

func (m *mergeMapSubscriberV2) NotifyNext(innerValue interface{}) {
	m.Destination.Next(innerValue)
}

func (m *mergeMapSubscriberV2) NotifyError(err error) {
	m.Destination.Error(err)
}

func (m *mergeMapSubscriberV2) NotifComplete(innerSub base.InnerSubscriber) {
	buffer := m.buffer
	m.Remove(innerSub)
	m.active--
	if len(buffer) > 0 {
		m.Next(buffer[0])
		buffer = buffer[1:]
	} else if m.active == 0 && m.hasCompleted {
		m.Destination.Complete()
	}
}

type mergeMapOperatorV2 struct {
	project    func(value interface{}) base.Subscribable
	concurrent int
}

func newMergeMapOperatorV2(project func(value interface{}) base.Subscribable, concurrent int) mergeMapOperatorV2 {
	newInstance := new(mergeMapOperatorV2)
	newInstance.project = project
	newInstance.concurrent = concurrent
	return *newInstance
}

func (m *mergeMapOperatorV2) Call(subscriber base.SubscriberLike, source base.Observable) base.SubscriptionLike {
	nmmsv2 := newMergeMapSubscriberV2(subscriber, m.project, m.concurrent)
	return source.Subscribe(&nmmsv2)
}

func MergeMapV2(project func(value interface{}) base.Subscribable) base.OperatorFunction {
	return func(source base.Observable) base.Observable {
		op := newMergeMapOperatorV2(project, 99999)
		ob := source.Lift(&op)
		return ob
	}
}
