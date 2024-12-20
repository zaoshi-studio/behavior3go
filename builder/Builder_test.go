package builder

import (
	b3 "github.com/magicsea/behavior3go"
	"github.com/magicsea/behavior3go/actions"
	"github.com/magicsea/behavior3go/composites"
	"github.com/magicsea/behavior3go/core"
	"github.com/magicsea/behavior3go/decorators"
	"github.com/magicsea/behavior3go/examples/share"
	"testing"
	"time"
)

func TestBuilderBase(t *testing.T) {

	_ = AdaptComposite(&composites.Sequence{})
	_ = AdaptDecorator(&decorators.Repeater{})
	_ = AdaptAction(&actions.Log{})

	builder := NewBuilder()

	//	首次添加节点就直接记录为根节点
	seqNode := builder.addChild(AdaptComposite(&composites.Sequence{}))
	{
		//	在 sequence 节点中添加一个 repeater 表示循环 2次
		repeaterNode := seqNode.addChild(AdaptDecorator(&decorators.Repeater{}, WithProperty("maxLoop", 2.0)))
		//	被 repeater 修饰的节点，执行 log action
		repeaterNode.addChild(AdaptAction(&actions.Log{}, WithProperty("info", "test")))
	}

	bevTree := builder.Build()

	blackBoard := core.NewBlackboard()
	for i := 0; i < 3; i++ {
		bevTree.Tick(i, blackBoard)
	}
}

func TestBuilder(t *testing.T) {
	builder := NewBuilder()

	extMap := b3.NewRegisterStructMaps()

	t.Run("test: action after action", func(t *testing.T) {
		builder.Reset()
		builder.AddAction(AdaptAction(&actions.Log{}, WithProperty("info", "log...11")))
		builder.SetExtMap(extMap)
		bevTree := builder.Build()

		board := core.NewBlackboard()

		for i := 0; i < 5; i++ {
			bevTree.Tick(i, board)
		}
	})

	t.Run("examples: load_from_tree", func(t *testing.T) {
		builder.Reset()
		memSequenceNode := builder.addChild(AdaptComposite(&composites.MemSequence{}))
		{
			repeater := memSequenceNode.addChild(AdaptDecorator(&decorators.Repeater{}, WithProperty("maxLoop", 2.0)))
			{
				repeater.addChild(AdaptAction(&actions.Log{}, WithProperty("info", "log...11")))
			}

			memSequenceNode.addChild(AdaptAction(&actions.Log{}, WithProperty("info", "log...22")))
			limiter := memSequenceNode.addChild(AdaptDecorator(&decorators.Limiter{}, WithProperty("maxLoop", 2.0)))
			{
				limiter.addChild(AdaptAction(&actions.Log{}, WithProperty("info", "log...333")))
			}
		}

		builder.SetExtMap(extMap)
		bevTree := builder.Build()

		board := core.NewBlackboard()

		for i := 0; i < 5; i++ {
			bevTree.Tick(i, board)
		}
	})

	t.Run("examples: test_work", func(t *testing.T) {
		builder.Reset()

		extMap.Register("SetValue", new(share.SetValue))
		extMap.Register("IsValue", new(share.IsValue))

		priority := builder.addChild(AdaptComposite(&composites.Priority{}))
		{
			memSequence1 := priority.addChild(AdaptComposite(&composites.MemSequence{}))
			{
				memSequence1.addChild(AdaptAction(&share.IsValue{},
					WithProperties(map[string]interface{}{
						"key":   "job",
						"value": 0.0,
					}),
					WithTitle("IsValue(<key>,<value>)1"),
				))
				memSequence1.addChild(AdaptAction(&actions.Log{},
					WithProperty("info", "job 0-0"),
					WithTitle("Log(<info>)"),
				))
				memSequence1.addChild(AdaptAction(&actions.Wait{},
					WithProperty("milliseconds", 1000.0),
					WithTitle("Wait0-1 <milliseconds>ms"),
				))
				memSequence1.addChild(AdaptAction(&actions.Log{},
					WithProperty("info", "job 0-1"),
					WithTitle("Log(<info>)"),
				))
				memSequence1.addChild(AdaptAction(&share.SetValue{},
					WithProperties(map[string]interface{}{
						"key":   "job",
						"value": 1.0,
					}),
					WithTitle("SetValue(<key>,<value>)1"),
				))
				memSequence1.addChild(AdaptAction(&actions.Log{},
					WithProperty("info", "job 0-2"),
					WithTitle("Log(<info>)"),
				))
				memSequence1.addChild(AdaptAction(&actions.Wait{},
					WithProperty("milliseconds", 100.0),
					WithTitle("Wait0-2 <milliseconds>ms"),
				))
			}

			memSequence2 := priority.addChild(AdaptComposite(&composites.MemSequence{}))
			{
				memSequence2.addChild(AdaptAction(&share.IsValue{},
					WithProperties(map[string]interface{}{
						"key":   "job",
						"value": 1.0,
					}),
					WithTitle("IsValue(<key>,<value>)2"),
				))
				memSequence2.addChild(AdaptAction(&actions.Log{},
					WithProperty("info", "job 1-0"),
					WithTitle("Log(<info>)"),
				))
				memSequence2.addChild(AdaptAction(&actions.Wait{},
					WithProperty("milliseconds", 1000.0),
					WithTitle("Wait1-1 <milliseconds>ms"),
				))
				memSequence2.addChild(AdaptAction(&actions.Log{},
					WithProperty("info", "job 1-1"),
					WithTitle("Log(<info>)"),
				))
				memSequence2.addChild(AdaptAction(&share.SetValue{},
					WithProperties(map[string]interface{}{
						"key":   "job",
						"value": 0.0,
					}),
					WithTitle("SetValue(<key>,<value>)2"),
				))
				memSequence2.addChild(AdaptAction(&actions.Log{},
					WithProperty("info", "job 1-2"),
					WithTitle("Log(<info>)"),
				))
				memSequence2.addChild(AdaptAction(&actions.Wait{},
					WithProperty("milliseconds", 100.0),
					WithTitle("Wait1-2 <milliseconds>ms"),
				))
			}
		}

		builder.SetExtMap(extMap)
		bevTree := builder.Build()
		bevTree.Print()

		board := core.NewBlackboard()
		for i := 0; i < 40; i++ {
			bevTree.Tick(i, board)
			time.Sleep(time.Millisecond * 100)
		}
	})

	t.Run("examples: subtree", func(t *testing.T) {
		builder.Reset()
		mem := builder.AddComposite(AdaptComposite(&composites.MemSequence{},
			WithTitle("MemSequence"),
		))
		{
			repeater := mem.AddDecorator(AdaptDecorator(&decorators.Repeater{},
				WithProperty("maxLoop", 2.0),
				WithTitle("Repeat <maxLoop>x"),
			))
			{
				repeater.AddAction(AdaptAction(&actions.Log{},
					WithProperty("info", "log...11"),
					WithTitle("Log"),
				))
			}

			mem.AddAction(AdaptAction(&actions.Log{}, WithProperty("info", "log...22")))
			mem.AddSubTree(AdaptAction(&actions.Log{}, WithProperty("info", " call child")))

			limiter := mem.AddDecorator(AdaptDecorator(&decorators.Limiter{}, WithProperty("maxLoop", 2.0)))
			{
				limiter.AddAction(AdaptAction(&actions.Log{}, WithProperty("info", "log...333")))
			}
		}

		bevTree := builder.Build()

		//tree2 := builder.AddTree(AdaptAction(&actions.Log{}, WithProperty("info", " call child")))
		board := core.NewBlackboard()
		for i := 0; i < 5; i++ {
			bevTree.Tick(i, board)
		}
	})
}
