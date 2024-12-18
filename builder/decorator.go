package builderv2

import (
	"github.com/magicsea/behavior3go/core"
	"reflect"
)

type IDecoratorAdapter interface {
	core.IDecorator
	IAdaptNode
	AdaptSetChild(adapter IAdaptNode)
}

func AdaptDecorator(decorator core.IDecorator, options ...Option) IDecoratorAdapter {
	adapter := &DecoratorAdapter{
		IDecorator: decorator,
		AdaptNode:  newAdaptNode(reflect.TypeOf(decorator).Elem().Name()),
	}
	for _, option := range options {
		option.apply(&adapter.AdaptNode)
	}
	return adapter
}

type DecoratorAdapter struct {
	core.IDecorator
	AdaptNode
}

func (adapter *DecoratorAdapter) GetID() string {
	return adapter.AdaptNode.GetID()
}
func (adapter *DecoratorAdapter) GetName() string {
	return adapter.AdaptNode.GetName()
}
func (adapter *DecoratorAdapter) GetTitle() string       { return adapter.AdaptNode.GetTitle() }
func (adapter *DecoratorAdapter) GetDescription() string { return adapter.AdaptNode.GetDescription() }

func (adapter *DecoratorAdapter) AdaptSetChild(adaptNode IAdaptNode) {
	adapter.IDecorator.SetChild(adaptNode)
}
