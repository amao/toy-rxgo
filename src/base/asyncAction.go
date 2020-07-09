package base

import (
	"time"
)

type AsyncAction struct {
	*Action
	id      interface{}
	state   interface{}
	delay   float64
	pending bool // default value = false
}

func NewAsyncAction(scheduler SchedulerLike, work func(SchedulerAction, interface{})) AsyncAction {
	newInstance := new(AsyncAction)
	action := NewAction(scheduler, work)
	newInstance.Action = &action
	return *newInstance
}

func (a *AsyncAction) Schedule(state interface{}, delay float64) SubscriptionLike {
	if a.closed {
		return a
	}

	a.state = state
	id := a.id
	scheduler := a.scheduler

	if id != nil {
		a.id = a.recycleAsyncId(scheduler, id, delay)
	}

	a.pending = true
	a.delay = delay

	if a.id == nil {
		a.requestAsyncId(scheduler, a.id, delay)
	}

	return a
}

func (a *AsyncAction) requestAsyncId(scheduler SchedulerLike, id interface{}, delay float64) interface{} {
	t := time.After(time.Duration(delay) * time.Millisecond)
	go func() {
		for range t {
			asyncScheduler := scheduler.(*AsyncScheduler)
			asyncScheduler.flush(*a)
		}
	}()

	return t
}

func (a *AsyncAction) recycleAsyncId(scheduler SchedulerLike, id interface{}, delay float64) interface{} {
	if a.delay == delay && a.pending == false {
		return id
	}

	ticker := id.(time.Ticker)
	ticker.Stop()

	return nil
}

func (a *AsyncAction) _execute(state interface{}, delay float64) interface{} {
	a.work(a, state)
	return nil
}

func (a *AsyncAction) execute(state interface{}, delay float64) interface{} {
	if a.closed {
		panic("executing a cancelled action")
	}

	a.pending = false

	a._execute(state, delay)

	if a.pending == false && a.id != nil {
		a.id = a.recycleAsyncId(a.scheduler, a.id, 0)
	}

	return nil
}
