package cvnode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffcv"
)

type StackProps struct {
	*scaffcv.AcceptChildren
}

func Stack(create func(t *scaff.Tracker, props *StackProps)) scaffcv.NodeBuilder {
	return scaffcv.Standard(scaffcv.StandardCreate[StackProps]{
		ID: "stack",
		DefaultProps: StackProps{
			AcceptChildren: &scaffcv.AcceptChildren{},
		},
		PropsCreator: create,
	})
}
