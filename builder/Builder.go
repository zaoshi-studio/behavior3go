package builder

import (
	b3 "github.com/magicsea/behavior3go"
	"github.com/magicsea/behavior3go/config"
	"github.com/magicsea/behavior3go/core"
	"github.com/magicsea/behavior3go/loader"
	"sync"
)

type nodeOperator struct {
	builder  *Builder
	parentID string
	treeID   string
}

func (op *nodeOperator) addChild(node IAdaptNode) *nodeOperator {
	op.builder.addNode(op.treeID, op.parentID, node)
	return &nodeOperator{
		builder:  op.builder,
		parentID: node.GetID(),
		treeID:   node.GetTreeID(),
	}
}

// AddAction
//
//	 @receiver c
//	 @param actNode
//	 @Description:
//		这里语义上定义只能添加 action 节点, 并且禁止链式调用
func (op *nodeOperator) AddAction(actNode IActionAdapter) {
	op.builder.addNode(op.treeID, op.parentID, actNode)
}

func (op *nodeOperator) AddComposite(compositeNode ICompositeAdapter) *nodeOperator {
	return op.addChild(compositeNode)
}

func (op *nodeOperator) AddDecorator(decoratorNode IDecoratorAdapter) *nodeOperator {
	return op.addChild(decoratorNode)
}

func (op *nodeOperator) AddCondition(conditionNode IConditionAdapter) *nodeOperator {
	return op.addChild(conditionNode)
}

func (op *nodeOperator) AddSubTree(node ISubTreeAdapter) {
	//	把当前节点作为一颗新树的根节点
	treeCfg := op.builder.addTree(node)

	//	同时以这个根节点的 id 作为 name 添加一个哑结点
	op.builder.addNode(op.treeID, op.parentID, &SubTreeAdapter{
		AdaptNode: newAdaptNode(treeCfg.ID, WithCategory("tree")),
	})
}

func (op *nodeOperator) AddTree(node IAdaptNode) *nodeOperator {
	op.builder.addTree(node)
	return &nodeOperator{
		builder:  op.builder,
		parentID: node.GetID(),
		treeID:   node.GetTreeID(),
	}
}

type Builder struct {
	nodeOperator

	subTreeLoadFunc func(id string) *core.BehaviorTree
	extMap          *b3.RegisterStructMaps
	treeMap         sync.Map
	projectCfg      *config.BTProjectCfg
	treeCfgs        map[string]*config.BTTreeCfg
	nodeCfgs        map[string]*config.BTNodeCfg
	adaptNodes      map[string]IAdaptNode
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

func (b *Builder) addTree(node IAdaptNode) *config.BTTreeCfg {

	nodeCfg := b.convertAdapterNodeToCfg(node)

	treeID := genTreeID()
	treeCfg := &config.BTTreeCfg{
		ID:          treeID,
		Title:       node.GetTitle(),
		Description: node.GetDescription(),
		Root:        node.GetID(),
		Properties:  make(map[string]interface{}),
		Nodes:       make(map[string]*config.BTNodeCfg),
	}
	for k, v := range node.GetProperties() {
		treeCfg.Properties[k] = v
	}

	b.treeCfgs[treeID] = treeCfg

	node.setTreeID(treeID)
	b.nodeCfgs[node.GetID()] = nodeCfg
	treeCfg.Nodes[node.GetID()] = nodeCfg
	b.adaptNodes[node.GetID()] = node

	//	当前没有树, 就当做第一个树, 并且默认勾选
	if len(b.projectCfg.Trees) == 0 {
		b.projectCfg.Select = treeCfg.ID
		b.nodeOperator.treeID = treeID
		b.nodeOperator.parentID = node.GetID()
	}

	b.projectCfg.Trees = append(b.projectCfg.Trees, treeCfg)

	return treeCfg
}

func (b *Builder) addNode(treeID, parentID string, node IAdaptNode) bool {
	if len(parentID) == 0 && len(b.projectCfg.Trees) != 0 {
		panic("parentID is empty")
	}

	if len(treeID) == 0 {
		b.addTree(node)
		return true
	}

	nodeCfg := b.convertAdapterNodeToCfg(node)

	treeCfg := b.treeCfgs[treeID]
	treeCfg.Nodes[node.GetID()] = nodeCfg

	node.setTreeID(treeID)

	b.nodeCfgs[node.GetID()] = nodeCfg
	b.adaptNodes[node.GetID()] = node

	parentNodeCfg := b.nodeCfgs[parentID]

	parentAdaptNode := b.adaptNodes[parentID]
	switch parentNodeCfg.Category {
	case b3.COMPOSITE:

		v, ok := parentAdaptNode.(ICompositeAdapter)
		if !ok {
			panic("wrong type")
		}

		v.AdaptAddChild(node)
		parentNodeCfg.Children = append(parentNodeCfg.Children, node.GetID())
	case b3.DECORATOR:
		v, ok := parentAdaptNode.(IDecoratorAdapter)
		if !ok {
			panic("wrong type")
		}
		v.AdaptSetChild(node)
		parentNodeCfg.Child = node.GetID()
	}

	//parentNodeCfg, ok := b.bevTreeCfg.Nodes[parentID]
	//if !ok && len(b.bevTreeCfg.Root) != 0 {
	//	return false
	//}
	//
	//_, ok = b.bevTreeCfg.Nodes[node.GetID()]
	//if ok {
	//	panic(fmt.Errorf("node: %v has been registered", node.GetID()))
	//}
	//
	//b.bevTreeCfg.Nodes[node.GetID()] = b.convertAdapterNodeToCfg(node)
	//b.adaptNodes[node.GetID()] = node
	//
	//if len(parentID) == 0 {
	//	b.bevTreeCfg.Root = node.GetID()
	//	b.nodeOperator.parentID = node.GetID()
	//	return true
	//} else {
	//	parentNodeCfg.Children = append(parentNodeCfg.Children, node.GetID())
	//}
	//
	//parentAdaptNode, ok := b.adaptNodes[parentID]
	//if !ok {
	//	panic(fmt.Sprintf("adapt node %v not found", node.GetID()))
	//}
	//
	//switch parentNodeCfg.Category {
	//case b3.COMPOSITE:
	//
	//	v, ok := parentAdaptNode.(ICompositeAdapter)
	//	if !ok {
	//		panic("wrong type")
	//	}
	//
	//	v.addChild(node)
	//case b3.DECORATOR:
	//	v, ok := parentAdaptNode.(IDecoratorAdapter)
	//	if !ok {
	//		panic("wrong type")
	//	}
	//	v.AdaptSetChild(node)
	//	parentNodeCfg.Child = node.GetID()
	//}

	return true
}

func (b *Builder) addNodeAsSubTree(treeID, parentID string) {

}

func (b *Builder) SetExtMap(extMap *b3.RegisterStructMaps) {
	b.extMap = extMap
}

func (b *Builder) Build() *core.BehaviorTree {
	var selectedTree *core.BehaviorTree
	//载入
	for _, v := range b.projectCfg.Trees {
		bevTree := loader.CreateBevTreeFromConfig(v, b.extMap)
		if selectedTree == nil {
			selectedTree = bevTree
		}
		b.treeMap.Store(v.ID, bevTree)
	}

	core.SetSubTreeLoadFunc(func(id string) *core.BehaviorTree {
		println("==>load subtree:", id)
		t, ok := b.treeMap.Load(id)
		if ok {
			return t.(*core.BehaviorTree)
		}
		return nil
	})
	return selectedTree
}

func (b *Builder) Reset() {
	b.nodeOperator = nodeOperator{}
	b.projectCfg = &config.BTProjectCfg{
		ID:    b3.CreateUUID(),
		Scope: "tree",
	}
	b.treeCfgs = make(map[string]*config.BTTreeCfg)
	b.nodeCfgs = make(map[string]*config.BTNodeCfg)
	b.adaptNodes = make(map[string]IAdaptNode)
	b.nodeOperator.builder = b
	b.extMap = b3.NewRegisterStructMaps()
	b.treeMap = sync.Map{}
}

func NewBuilder() *Builder {
	builder := &Builder{
		nodeOperator: nodeOperator{},
		projectCfg: &config.BTProjectCfg{
			ID:    b3.CreateUUID(),
			Scope: "tree",
		},
		treeCfgs:   make(map[string]*config.BTTreeCfg),
		nodeCfgs:   make(map[string]*config.BTNodeCfg),
		adaptNodes: make(map[string]IAdaptNode),
	}
	builder.nodeOperator.builder = builder

	return builder
}

func genNodeID() string {
	//return fmt.Sprintf("%v", time.Now().Nanosecond())
	return b3.CreateUUID()
}

func genTreeID() string {
	return b3.CreateUUID()
}
