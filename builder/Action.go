package builder

import (
	"github.com/magicsea/behavior3go/core"
	"reflect"
)

type IActionAdapter interface {
	core.IAction
	IAdaptNode
}

func AdaptAction(action core.IAction, options ...Option) IActionAdapter {
	adapter := &ActionAdapter{
		IAction:   action,
		AdaptNode: newAdaptNode(reflect.TypeOf(action).Elem().Name(), options...),
	}
	return adapter
}

type ActionAdapter struct {
	core.IAction
	AdaptNode
}

func (adapter *ActionAdapter) GetID() string {
	return adapter.AdaptNode.GetID()
}
func (adapter *ActionAdapter) GetName() string {
	return adapter.AdaptNode.GetName()
}
func (adapter *ActionAdapter) GetTitle() string       { return adapter.AdaptNode.GetTitle() }
func (adapter *ActionAdapter) GetCategory() string    { return adapter.IAction.GetCategory() }
func (adapter *ActionAdapter) GetDescription() string { return adapter.AdaptNode.GetDescription() }
