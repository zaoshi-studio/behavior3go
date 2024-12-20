package builder

import "github.com/magicsea/behavior3go/core"

type IAdaptNode interface {
	core.IBaseNode
	GetID() string
	setTreeID(string)
	GetTreeID() string
	GetCategory() string
	GetParameters() map[string]interface{}
	GetProperties() map[string]interface{}
}

type AdaptNode struct {
	id          string
	name        string
	title       string
	treeID      string
	category    string
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
func (node *AdaptNode) GetTitle() string        { return node.title }
func (node *AdaptNode) setTreeID(treeID string) { node.treeID = treeID }
func (node *AdaptNode) GetTreeID() string       { return node.treeID }
func (node *AdaptNode) GetCategory() string     { return node.category }
func (node *AdaptNode) GetDescription() string  { return node.description }
func (node *AdaptNode) GetParameters() map[string]interface{} {
	return node.parameters
}
func (node *AdaptNode) GetProperties() map[string]interface{} {
	return node.properties
}

func newAdaptNode(name string, options ...Option) AdaptNode {
	node := AdaptNode{
		id:         genNodeID(),
		name:       name,
		properties: make(map[string]interface{}),
		parameters: make(map[string]interface{}),
	}

	for _, option := range options {
		option.apply(&node)
	}
	return node
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

func WithCategory(category string) Option {
	return func(adapter *AdaptNode) {
		adapter.category = category
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
