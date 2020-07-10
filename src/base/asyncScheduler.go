package base

import (
	"reflect"
	"time"
)

type AsyncScheduler struct {
	*Scheduler
	delegate  *Scheduler
	actions   []AsyncAction
	active    bool //default value = false
	scheduled interface{}
}

func NewAsyncScheduler(schedulerAction string, now func() time.Time) AsyncScheduler {
	newInstance := new(AsyncScheduler)
	parentScheduler := NewScheduler(schedulerAction, func() time.Time {
		if newInstance.delegate != nil && !reflect.DeepEqual(newInstance.delegate, newInstance) {
			return newInstance.delegate.Now()
		} else {
			return now()
		}
	})
	newInstance.Scheduler = &parentScheduler
	return *newInstance
}

func (a *AsyncScheduler) Schedule(scheduler SchedulerLike, work func(SchedulerAction, interface{}), delay uint, state interface{}) SubscriptionLike {
	if a.delegate != nil && !reflect.DeepEqual(a.delegate, a) {
		return a.delegate.Schedule(a, work, delay, state)
	} else {
		return a.Scheduler.Schedule(a, work, delay, state)
	}
}

func (a *AsyncScheduler) flush(action *AsyncAction) {
	actions := a.actions

	if a.active {
		actions = append(actions, *action)
		return
	}

	a.active = true

	action.execute(action.state, action.delay)

	for _, action := range actions {
		action.execute(action.state, action.delay)
	}

	a.active = false
}
