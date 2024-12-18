package builderv2

import (
	b3 "github.com/magicsea/behavior3go"
	"github.com/magicsea/behavior3go/actions"
	"github.com/magicsea/behavior3go/composites"
	"github.com/magicsea/behavior3go/core"
	"github.com/magicsea/behavior3go/decorators"
	"github.com/magicsea/behavior3go/examples/share"
	"github.com/magicsea/behavior3go/loader"
	"sync"
	"testing"
	"time"
)

func TestBuilderBase(t *testing.T) {

	_ = AdaptComposite(&composites.Sequence{})
	_ = AdaptDecorator(&decorators.Repeater{})
	_ = AdaptAction(&actions.Log{})

	builder := NewBuilder()

	//	首次添加节点就直接记录为根节点
	seqNode := builder.AddChild(AdaptComposite(&composites.Sequence{}))
	{
		//	在 sequence 节点中添加一个 repeater 表示循环 2次
		repeaterNode := seqNode.AddChild(AdaptDecorator(&decorators.Repeater{}, WithProperty("maxLoop", 2.0)))
		//	被 repeater 修饰的节点，执行 log action
		repeaterNode.AddChild(AdaptAction(&actions.Log{}, WithProperty("info", "test")))
	}

	bevTree := loader.CreateBevTreeFromConfig(builder.BevTreeCfg(), nil)

	blackBoard := core.NewBlackboard()
	for i := 0; i < 3; i++ {
		bevTree.Tick(i, blackBoard)
	}
}

func TestBuilder(t *testing.T) {
	builder := NewBuilder()

	extMap := b3.NewRegisterStructMaps()

	t.Run("examples: load_from_tree", func(t *testing.T) {
		memSequenceNode := builder.AddChild(AdaptComposite(&composites.MemSequence{}))
		{
			repeater := memSequenceNode.AddChild(AdaptDecorator(&decorators.Repeater{}, WithProperty("maxLoop", 2.0)))
			{
				repeater.AddChild(AdaptAction(&actions.Log{}, WithProperty("info", "log...11")))
			}
			memSequenceNode.AddChild(AdaptAction(&actions.Log{}, WithProperty("info", "log...22")))
			limiter := memSequenceNode.AddChild(AdaptDecorator(&decorators.Limiter{}, WithProperty("maxLoop", 2.0)))
			{
				limiter.AddChild(AdaptAction(&actions.Log{}, WithProperty("info", "log...333")))
			}
		}

		bevTree := loader.CreateBevTreeFromConfig(builder.BevTreeCfg(), extMap)

		board := core.NewBlackboard()

		for i := 0; i < 5; i++ {
			bevTree.Tick(i, board)
		}
	})

	t.Run("examples: test_work", func(t *testing.T) {
		var mapTreesByID = sync.Map{}

		extMap.Register("SetValue", new(share.SetValue))
		extMap.Register("IsValue", new(share.IsValue))

		priority := builder.AddChild(AdaptComposite(&composites.Priority{}))
		{
			memSequence1 := priority.AddChild(AdaptComposite(&composites.MemSequence{}))
			{
				memSequence1.AddChild(AdaptAction(&share.IsValue{},
					WithProperties(map[string]interface{}{
						"key":   "job",
						"value": 0.0,
					}),
					WithTitle("IsValue(<key>,<value>)1"),
				))
				memSequence1.AddChild(AdaptAction(&actions.Log{},
					WithProperty("info", "job 0-0"),
					WithTitle("Log(<info>)"),
				))
				memSequence1.AddChild(AdaptAction(&actions.Wait{},
					WithProperty("milliseconds", 1000.0),
					WithTitle("Wait0-1 <milliseconds>ms"),
				))
				memSequence1.AddChild(AdaptAction(&actions.Log{},
					WithProperty("info", "job 0-1"),
					WithTitle("Log(<info>)"),
				))
				memSequence1.AddChild(AdaptAction(&share.SetValue{},
					WithProperties(map[string]interface{}{
						"key":   "job",
						"value": 1.0,
					}),
					WithTitle("SetValue(<key>,<value>)1"),
				))
				memSequence1.AddChild(AdaptAction(&actions.Log{},
					WithProperty("info", "job 0-2"),
					WithTitle("Log(<info>)"),
				))
				memSequence1.AddChild(AdaptAction(&actions.Wait{},
					WithProperty("milliseconds", 100.0),
					WithTitle("Wait0-2 <milliseconds>ms"),
				))
			}

			memSequence2 := priority.AddChild(AdaptComposite(&composites.MemSequence{}))
			{
				memSequence2.AddChild(AdaptAction(&share.IsValue{},
					WithProperties(map[string]interface{}{
						"key":   "job",
						"value": 1.0,
					}),
					WithTitle("IsValue(<key>,<value>)2"),
				))
				memSequence2.AddChild(AdaptAction(&actions.Log{},
					WithProperty("info", "job 1-0"),
					WithTitle("Log(<info>)"),
				))
				memSequence2.AddChild(AdaptAction(&actions.Wait{},
					WithProperty("milliseconds", 1000.0),
					WithTitle("Wait1-1 <milliseconds>ms"),
				))
				memSequence2.AddChild(AdaptAction(&actions.Log{},
					WithProperty("info", "job 1-1"),
					WithTitle("Log(<info>)"),
				))
				memSequence2.AddChild(AdaptAction(&share.SetValue{},
					WithProperties(map[string]interface{}{
						"key":   "job",
						"value": 0.0,
					}),
					WithTitle("SetValue(<key>,<value>)2"),
				))
				memSequence2.AddChild(AdaptAction(&actions.Log{},
					WithProperty("info", "job 1-2"),
					WithTitle("Log(<info>)"),
				))
				memSequence2.AddChild(AdaptAction(&actions.Wait{},
					WithProperty("milliseconds", 100.0),
					WithTitle("Wait1-2 <milliseconds>ms"),
				))
			}
		}

		bevTree := loader.CreateBevTreeFromConfig(builder.BevTreeCfg(), extMap)
		bevTree.Print()

		mapTreesByID.Store(bevTree.GetID(), bevTree)
		core.SetSubTreeLoadFunc(func(id string) *core.BehaviorTree {
			tree, ok := mapTreesByID.Load(id)
			if ok {
				return tree.(*core.BehaviorTree)
			}
			return nil
		})

		board := core.NewBlackboard()
		for i := 0; i < 40; i++ {
			bevTree.Tick(i, board)
			time.Sleep(time.Millisecond * 100)
		}
	})
}
