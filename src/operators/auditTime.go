package operators

import (
	"github.com/amao/toy-rxgo/src/observables"

	"github.com/amao/toy-rxgo/src/base"
)

func AuditTime(duration int) base.OperatorFunction {
	return Audit(func(value interface{}) base.Subscribable {
		return observables.Timer(duration, duration)
	})
}
