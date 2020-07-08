package base

type Action struct {
	*Subscription
	scheduler SchedulerLike
	work      func(state interface{})
}

func NewAction(scheduler SchedulerLike, work func(interface{})) Action {
	newInstance := new(Action)
	parentSubscription := NewSubscription(nil)
	newInstance.Subscription = &parentSubscription
	newInstance.scheduler = scheduler
	newInstance.work = work
	return *newInstance
}

func (a *Action) Schedule(state interface{}, delay float64) SubscriptionLike {
	return a
}
