package base

type EmptyObserver struct {
	closed   bool
	next     func(interface{})
	error_   func(error)
	complete func()
}

func NewEmptyObserver() EmptyObserver {
	baseSubscriber := new(EmptyObserver)
	baseSubscriber.closed = true
	baseSubscriber.next = func(value interface{}) {}
	baseSubscriber.error_ = func(e error) {}
	baseSubscriber.complete = func() {}
	return *baseSubscriber
}

func (eo *EmptyObserver) Next(value interface{}) {
	eo.next(value)
}

func (eo *EmptyObserver) Error(e error) {
	eo.error_(e)
}

func (eo *EmptyObserver) Complete() {
	eo.complete()
}
