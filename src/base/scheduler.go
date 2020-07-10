package base

import "time"

type Scheduler struct {
	now             func() time.Time
	schedulerAction string
}

func NewScheduler(schedulerAction string, now func() time.Time) Scheduler {
	newInstance := new(Scheduler)
	newInstance.now = now
	newInstance.schedulerAction = schedulerAction
	return *newInstance
}

func (s *Scheduler) Now() time.Time {
	return s.now()
}

func (s *Scheduler) Schedule(scheduler SchedulerLike, work func(SchedulerAction, interface{}), delay uint, state interface{}) SubscriptionLike {
	switch s.schedulerAction {
	case "AsyncAction":
		action := NewAsyncAction(scheduler, work)
		return action.Schedule(state, delay)
	default:
		action := NewAction(s, work)
		return action.Schedule(state, delay)
	}
}
