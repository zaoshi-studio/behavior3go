package builder

import (
	"github.com/magicsea/behavior3go/core"
	"reflect"
)

type IConditionAdapter interface {
	core.ICondition
	IAdaptNode
}

func AdaptCondition(condition core.ICondition, options ...Option) IConditionAdapter {
	adapter := &ConditionAdapter{
		ICondition: condition,
		AdaptNode:  newAdaptNode(reflect.TypeOf(condition).Elem().Name(), options...),
	}
	return adapter
}

type ConditionAdapter struct {
	core.ICondition
	AdaptNode
}

func (adapter *ConditionAdapter) GetID() string {
	return adapter.AdaptNode.GetID()
}
func (adapter *ConditionAdapter) GetName() string {
	return adapter.AdaptNode.GetName()
}
func (adapter *ConditionAdapter) GetTitle() string       { return adapter.AdaptNode.GetTitle() }
func (adapter *ConditionAdapter) GetCategory() string    { return adapter.ICondition.GetCategory() }
func (adapter *ConditionAdapter) GetDescription() string { return adapter.AdaptNode.GetDescription() }