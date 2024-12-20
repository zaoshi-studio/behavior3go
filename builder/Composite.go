package builder

import (
	"github.com/magicsea/behavior3go/core"
	"reflect"
)

type ICompositeAdapter interface {
	core.IComposite
	IAdaptNode
	AdaptAddChild(child IAdaptNode)
}

func AdaptComposite(composite core.IComposite, options ...Option) ICompositeAdapter {
	adapter := &CompositeAdapter{
		IComposite: composite,
		AdaptNode:  newAdaptNode(reflect.TypeOf(composite).Elem().Name(), options...),
	}
	return adapter
}

type CompositeAdapter struct {
	core.IComposite
	AdaptNode
}

func (adapter *CompositeAdapter) GetID() string          { return adapter.AdaptNode.GetID() }
func (adapter *CompositeAdapter) GetName() string        { return adapter.AdaptNode.GetName() }
func (adapter *CompositeAdapter) GetTitle() string       { return adapter.AdaptNode.GetTitle() }
func (adapter *CompositeAdapter) GetCategory() string    { return adapter.IComposite.GetCategory() }
func (adapter *CompositeAdapter) GetDescription() string { return adapter.AdaptNode.GetDescription() }

func (adapter *CompositeAdapter) AdaptAddChild(child IAdaptNode) {
	adapter.IComposite.AddChild(child)
}
