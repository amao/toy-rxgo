package operators

import (
	"math"

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
	self := base.NewSubscriber(destination)
	super := new(base.OuterSubscriber)
	newInstance.OuterSubscriber = super
	newInstance.Subscriber = &self
	newInstance.Destination.Add(newInstance)
	newInstance.project = project
	newInstance.concurrent = concurrent
	newInstance.SetInnerNext(func(value interface{}) {
		if newInstance.active < newInstance.concurrent {
			newInstance._tryNext(value)
		} else {
			newInstance.buffer = append(newInstance.buffer, value)
		}
	})
	newInstance.SetInnerComplete(func() {
		newInstance.hasCompleted = true
		if newInstance.active == 0 && len(newInstance.buffer) == 0 {
			newInstance.Destination.Complete()
		}
		newInstance.Unsubscribe()
	})
	return *newInstance
}

func (m *mergeMapSubscriberV2) _tryNext(value interface{}) {
	var result base.Subscribable
	result = m.project(value)
	m.active++
	m._innerSub(result, value)
}

func (m *mergeMapSubscriberV2) _innerSub(ish base.Subscribable, value interface{}) {
	innerSubscriber := base.NewInnerSubscriber(m, -1, -1)
	destination := m.Destination
	destination.Add(innerSubscriber)

	if innerSubscriber.Closed() {
		return
	}

	ish.Subscribe(&innerSubscriber)
}

func (m *mergeMapSubscriberV2) NotifyNext(outerValue interface{}, innerValue interface{}, outerIndex interface{}, innerIndex interface{}, innerSub base.InnerSubscriber) {
	m.Destination.Next(innerValue)
}

func (m *mergeMapSubscriberV2) NotifyError(err error) {
	m.Destination.Error(err)
}

func (m *mergeMapSubscriberV2) NotifyComplete(innerSub base.InnerSubscriber) {
	m.Remove(innerSub)
	m.active--
	if len(m.buffer) > 0 {
		m.CallInnerNext(m.buffer[0])
		m.buffer = m.buffer[1:]
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

func (m *mergeMapOperatorV2) Call(subscriber base.SubscriberLike, source base.Subscribable) base.SubscriptionLike {
	nmmsv2 := newMergeMapSubscriberV2(subscriber, m.project, m.concurrent)
	return source.Subscribe(&nmmsv2)
}

func MergeMapV2(project func(value interface{}) base.Subscribable, concurrent ...interface{}) base.OperatorFunction {
	concurrentN := math.MaxUint32
	if len(concurrent) == 1 {
		if value, ok := concurrent[0].(int); ok {
			concurrentN = value
		}
	}

	return func(source base.Observable) base.Observable {
		op := newMergeMapOperatorV2(project, concurrentN)
		ob := source.Lift(&op)
		return ob
	}
}
