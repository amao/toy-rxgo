package util

import "github.com/amao/toy-rxgo/src/base"

func Lift(source base.Observable, operator base.Operator) base.Observable {
	if HasLift(source) {
		return source.Lift(operator)
	}
	panic("Unable to lift unknown Observable type")
}

func HasLift(source interface{}) bool {
	ob := source.(base.Observable)
	return source != nil && ob.Lift != nil
}

func StankyLift(source base.Observable, liftedSource base.Observable, operator base.Operator) base.Observable {
	if HasLift(source) {
		return liftedSource.Lift(operator)
	}
	panic("Unable to lift unknown Observable type")
}
