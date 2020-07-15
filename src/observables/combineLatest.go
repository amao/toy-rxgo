package observables

import (
	"github.com/amao/toy-rxgo/src/base"
	"github.com/amao/toy-rxgo/src/util"
)

type combineLatestSubscriber struct {
	*base.OuterSubscriber
	active         int
	values         []interface{}
	observables    []interface{}
	toRespond      int
	resultSelector func(values []interface{}) interface{}
}

func newCombineLatestSubscriber(destination base.SubscriberLike, resultSelector func(values []interface{}) interface{}) combineLatestSubscriber {
	newInstance := new(combineLatestSubscriber)

	selfSubscriber := base.NewSubscriber(destination)
	selfOuterSubscriber := new(base.OuterSubscriber)
	newInstance.OuterSubscriber = selfOuterSubscriber
	newInstance.Subscriber = &selfSubscriber
	newInstance.Destination.Add(newInstance)

	newInstance.resultSelector = resultSelector

	newInstance.SetInnerNext(func(observable interface{}) {
		newInstance.values = append(newInstance.values, nil)
		newInstance.observables = append(newInstance.observables, observable)
	})

	newInstance.SetInnerComplete(func() {
		length := len(newInstance.observables)
		if length == 0 {
			newInstance.Destination.Complete()
		} else {
			newInstance.active = length
			newInstance.toRespond = length
			for index, observable := range newInstance.observables {
				innerSubscriber := base.NewInnerSubscriber(newInstance, observable, index)
				newInstance.Add(observable.(*base.Observable).Subscribe(innerSubscriber))
			}
		}
	})

	return *newInstance
}

func (c *combineLatestSubscriber) NotifyNext(outerValue interface{}, innerValue interface{}, outerIndex interface{}, innerIndex interface{}, innerSub base.InnerSubscriber) {
	values := c.values
	oldVal := values[outerIndex.(int)]
	var toRespond int
	if c.toRespond == 0 {
		toRespond = 0
	} else if oldVal == nil {
		c.toRespond -= 1
		toRespond = c.toRespond
	} else {
		toRespond = c.toRespond
	}

	values[outerIndex.(int)] = innerValue

	if toRespond == 0 {
		if c.resultSelector != nil {
			c._tryResultSelector(values)
		} else {
			c.Destination.Next(values[:])
		}
	}
}

func (c *combineLatestSubscriber) NotifyComplete(innerSub base.InnerSubscriber) {
	c.active -= 1
	if c.active == 0 {
		c.Destination.Complete()
	}
}

func (c *combineLatestSubscriber) _tryResultSelector(values []interface{}) {
	var result interface{}
	result = c.resultSelector(values)
	c.Destination.Next(result)
}

type combineLatestOperator struct {
	resultSelector func(values []interface{}) interface{}
}

func newCombineLatestOperator(resultSelector func(values []interface{}) interface{}) combineLatestOperator {
	newInstance := new(combineLatestOperator)
	newInstance.resultSelector = resultSelector
	return *newInstance
}

func (c *combineLatestOperator) Call(subscriber base.SubscriberLike, source base.Subscribable) base.SubscriptionLike {
	ncls := newCombineLatestSubscriber(subscriber, c.resultSelector)
	return source.Subscribe(&ncls)
}

func CombineLatest(observables ...interface{}) *base.Observable {
	obj := observables[len(observables)-1]
	var operator combineLatestOperator
	if resultSelector, ok := obj.(func([]interface{}) interface{}); ok {
		operator = newCombineLatestOperator(resultSelector)
		observables = observables[:len(observables)-1]
	} else {
		operator = *new(combineLatestOperator)
	}
	ob := util.Lift(FromArray(observables), &operator)
	return &ob
}
