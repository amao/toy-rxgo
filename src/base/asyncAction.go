package base

import (
	"time"
)

type AsyncAction struct {
	*Action
	id      interface{}
	state   interface{}
	delay   uint
	pending bool // default value = false
}

func NewAsyncAction(scheduler SchedulerLike, work func(SchedulerAction, interface{})) AsyncAction {
	newInstance := new(AsyncAction)
	action := NewAction(scheduler, work)
	newInstance.Action = &action
	return *newInstance
}

func (a *AsyncAction) Schedule(state interface{}, delay uint) SubscriptionLike {
	if a.closed {
		return a
	}

	a.state = state
	id := a.id
	scheduler := a.scheduler

	if _, ok := id.(struct{}); ok {
		a.id = a.recycleAsyncId(scheduler, id, delay)
	}

	//if id != nil {
	//	a.id = a.recycleAsyncId(scheduler, id, delay)
	//}

	a.pending = true
	a.delay = delay

	if a.id == nil {
		a.id = a.requestAsyncId(scheduler, a.id, delay)
	}

	return a
}

func (a *AsyncAction) requestAsyncId(scheduler SchedulerLike, id interface{}, delay uint) interface{} {
	// Why not using time.After? See implementation of AsyncAction in RxJs
	asyncScheduler := scheduler.(*AsyncScheduler)
	t := time.NewTicker(time.Duration(delay) * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-t.C:
				asyncScheduler.flush(a)
			}

		}
	}()

	return struct {
		t    *time.Ticker
		done chan bool
	}{t, done}
}

func (a *AsyncAction) recycleAsyncId(scheduler SchedulerLike, id interface{}, delay uint) interface{} {
	if a.delay == delay && a.pending == false {
		return id
	}

	td := id.(struct {
		t    *time.Ticker
		done chan bool
	})

	td.t.Stop()
	td.done <- true

	return nil
}

func (a *AsyncAction) _execute(state interface{}, delay uint) interface{} {
	a.work(a, state)
	return nil
}

func (a *AsyncAction) execute(state interface{}, delay uint) interface{} {
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

func (a *AsyncAction) Unsubscribe() {
	if a.id != nil {
		a.id = a.recycleAsyncId(a.scheduler, a.id, 0)
	}
}
