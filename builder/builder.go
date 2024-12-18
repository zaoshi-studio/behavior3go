package builderv2

import (
	"fmt"
	b3 "github.com/magicsea/behavior3go"
	"github.com/magicsea/behavior3go/config"
	"github.com/magicsea/behavior3go/core"
	"time"
)

type IAdaptNode interface {
	core.IBaseNode
	GetID() string
	GetParameters() map[string]interface{}
	GetProperties() map[string]interface{}
}

type AdaptNode struct {
	id          string
	name        string
	title       string
	description string
	properties  map[string]interface{}
	parameters  map[string]interface{}
}

func (node *AdaptNode) GetID() string {
	return node.id
}
func (node *AdaptNode) GetName() string {
	return node.name
}
func (node *AdaptNode) GetTitle() string       { return node.title }
func (node *AdaptNode) GetDescription() string { return node.description }
func (node *AdaptNode) GetParameters() map[string]interface{} {
	return node.parameters
}
func (node *AdaptNode) GetProperties() map[string]interface{} {
	return node.properties
}

func newAdaptNode(name string) AdaptNode {
	return AdaptNode{
		id:         genNodeID(),
		name:       name,
		properties: make(map[string]interface{}),
		parameters: make(map[string]interface{}),
	}
}

type Option func(adapter *AdaptNode)

func (f Option) apply(adapter *AdaptNode) {
	f(adapter)
}

func WithTitle(title string) Option {
	return func(adapter *AdaptNode) {
		adapter.title = title
	}
}

func WithDescription(description string) Option {
	return func(adapter *AdaptNode) {
		adapter.description = description
	}
}

func WithParameter(k string, v interface{}) Option {
	return func(adapter *AdaptNode) {
		adapter.parameters[k] = v
	}
}

func WithParameters(parameters map[string]interface{}) Option {
	return func(adapter *AdaptNode) {
		for k, v := range parameters {
			adapter.parameters[k] = v
		}
	}
}

func WithProperty(k string, v interface{}) Option {
	return func(adapter *AdaptNode) {
		adapter.properties[k] = v
	}
}

func WithProperties(properties map[string]interface{}) Option {
	return func(adapter *AdaptNode) {
		for k, v := range properties {
			adapter.properties[k] = v
		}
	}
}

type NodeController struct {
	builder *Builder
	ownerID string
}

func (c *NodeController) AddChild(node IAdaptNode) *NodeController {
	if !c.builder.addNodes(c.ownerID, node) {
		return nil
	}
	return &NodeController{
		builder: c.builder,
		ownerID: node.GetID(),
	}
}

type Builder struct {
	NodeController
	bevTreeCfg *config.BTTreeCfg
	adaptNodes map[string]IAdaptNode
}

func (b *Builder) convertAdapterNodeToCfg(node IAdaptNode) *config.BTNodeCfg {
	node.Ctor()
	cfg := &config.BTNodeCfg{
		Id:          node.GetID(),
		Name:        node.GetName(),
		Title:       node.GetTitle(),
		Category:    node.GetCategory(),
		Description: node.GetDescription(),
		Parameters:  make(map[string]interface{}),
		Properties:  make(map[string]interface{}),
	}

	for k, v := range node.GetProperties() {
		cfg.Properties[k] = v
	}

	for k, v := range node.GetParameters() {
		cfg.Parameters[k] = v
	}

	return cfg
}

func (b *Builder) addNodes(parentID string, node IAdaptNode) bool {
	parentNode, ok := b.bevTreeCfg.Nodes[parentID]
	if !ok && len(b.bevTreeCfg.Root) != 0 {
		return false
	}

	//	这里直接添加为根节点

	_, ok = b.bevTreeCfg.Nodes[node.GetID()]
	if ok {
		panic(fmt.Errorf("node: %v has been registered", node.GetID()))
	}

	b.bevTreeCfg.Nodes[node.GetID()] = b.convertAdapterNodeToCfg(node)
	b.adaptNodes[node.GetID()] = node

	if len(parentID) == 0 {
		b.bevTreeCfg.Root = node.GetID()
		b.NodeController.ownerID = node.GetID()
		return true
	} else {
		parentNode.Children = append(parentNode.Children, node.GetID())
	}

	adaptNode, ok := b.adaptNodes[parentID]
	if !ok {
		panic(fmt.Sprintf("adapt node %v not found", node.GetID()))
	}

	switch parentNode.Category {
	case b3.COMPOSITE:

		v, ok := adaptNode.(ICompositeAdapter)
		if !ok {
			panic("wrong type")
		}

		v.AddChild(node)
	case b3.DECORATOR:
		v, ok := adaptNode.(IDecoratorAdapter)
		if !ok {
			panic("wrong type")
		}
		v.AdaptSetChild(node)
		parentNode.Child = node.GetID()
	}

	return true
}

func (b *Builder) BevTreeCfg() *config.BTTreeCfg {
	return b.bevTreeCfg
}

func NewBuilder() *Builder {
	builder := &Builder{
		NodeController: NodeController{},
		bevTreeCfg: &config.BTTreeCfg{
			Nodes:      make(map[string]*config.BTNodeCfg),
			Properties: make(map[string]interface{}),
		},
		adaptNodes: make(map[string]IAdaptNode),
	}
	builder.NodeController.builder = builder
	return builder
}

func genNodeID() string {
	return fmt.Sprintf("%v", time.Now().Nanosecond())
}
