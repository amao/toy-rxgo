package base

type Action struct {
	*Subscription
	scheduler SchedulerLike
	work      func(SchedulerAction, interface{})
}

func NewAction(scheduler SchedulerLike, work func(SchedulerAction, interface{})) Action {
	newInstance := new(Action)
	parentSubscription := NewSubscription(nil)
	newInstance.Subscription = &parentSubscription
	newInstance.scheduler = scheduler
	newInstance.work = work
	return *newInstance
}

func (a *Action) Schedule(state interface{}, delay uint) SubscriptionLike {
	return a
}
