package uinode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffui"
)

type StackProps struct {
	*scaffui.AcceptChildren
}

func Stack(create func(t *scaff.Tracker, props *StackProps)) scaffui.NodeBuilder {

	// The default behavior of multi-node covers all of the things we want the stack to do, so no work for us :D
	return scaffui.MultiNode(scaffui.MultiNodeCreate[StackProps]{
		ID: "stack",
		DefaultProps: StackProps{
			AcceptChildren: &scaffui.AcceptChildren{},
		},
		PropsCreator: create,
	})
}
