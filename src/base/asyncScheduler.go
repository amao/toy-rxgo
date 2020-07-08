package base

import "reflect"

type AsyncScheduler struct {
	*Scheduler
	delegate  *Scheduler
	actions   []AsyncAction
	active    bool //default value = false
	scheduled interface{}
}

func NewAsyncScheduler(schedulerAction SchedulerAction, now func() float64) AsyncScheduler {
	newInstance := new(AsyncScheduler)
	parentScheduler := NewScheduler(schedulerAction, func() float64 {
		if newInstance.delegate != nil && !reflect.DeepEqual(newInstance.delegate, newInstance) {
			return newInstance.delegate.Now()
		} else {
			return now()
		}
	})
	newInstance.Scheduler = &parentScheduler
	return *newInstance
}

func (a *AsyncScheduler) Schedule(work func(interface{}), delay float64, state interface{}) SubscriptionLike {
	if a.delegate != nil && !reflect.DeepEqual(a.delegate, a) {
		return a.delegate.Schedule(work, delay, state)
	} else {
		return a.Scheduler.Schedule(work, delay, state)
	}
}

func (a *AsyncScheduler) flush(action AsyncAction) {
	actions := a.actions

	if a.active {
		actions = append(actions, action)
		return
	}

	a.active = true

	for _, action := range actions {
		action.execute(action.state, action.delay)
	}

	a.active = false
}
