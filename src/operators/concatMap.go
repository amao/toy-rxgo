package operators

import "github.com/amao/toy-rxgo/src/base"

func ConcatMap(project func(value interface{}) base.Subscribable) base.OperatorFunction {
	return MergeMapV2(project, 1)
}
