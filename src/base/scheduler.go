package base

type Scheduler struct {
	now             func() float64
	schedulerAction SchedulerAction
}

func NewScheduler(schedulerAction SchedulerAction, now func() float64) Scheduler {
	newInstance := new(Scheduler)
	newInstance.now = now
	newInstance.schedulerAction = schedulerAction
	return *newInstance
}

func (s *Scheduler) Now() float64 {
	return s.now()
}

func (s *Scheduler) Schedule(work func(interface{}), delay float64, state interface{}) SubscriptionLike {
	switch s.schedulerAction.(type) {
	case *AsyncAction:
		action := NewAsyncAction(s, work)
		return action.Schedule(state, delay)
	default:
		action := NewAction(s, work)
		return action.Schedule(state, delay)
	}
}
