package basenode

import (
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
)

type RootProps struct {
	child optional.O[scaffui.NodeBuilder]
}

func (rp *RootProps) Child(builder scaffui.NodeBuilder) {
	rp.child.SetValue(builder)
}

func Root(create func(t *scaffui.Tracker, props *RootProps)) scaffui.NodeBuilder {
	return scaffui.UseSingleNode("root", create, func(core *scaffui.SingleChildConstruct[RootProps]) {

		// Pass the child to the core node
		if child, ok := core.Props().child.Value(); ok {
			core.Child(child)
		}
	})
}
