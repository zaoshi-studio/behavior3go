package builder

import "github.com/magicsea/behavior3go/core"

type ISubTreeAdapter interface {
	core.IAction
	IAdaptNode
}

type SubTreeAdapter struct {
	core.Action
	AdaptNode
}

func (adapter *SubTreeAdapter) GetID() string {
	return adapter.AdaptNode.GetID()
}
func (adapter *SubTreeAdapter) GetName() string {
	return adapter.AdaptNode.GetName()
}
func (adapter *SubTreeAdapter) GetTitle() string       { return adapter.AdaptNode.GetTitle() }
func (adapter *SubTreeAdapter) GetCategory() string    { return "tree" }
func (adapter *SubTreeAdapter) GetDescription() string { return adapter.AdaptNode.GetDescription() }
