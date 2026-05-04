package basenode

import "github.com/Liphium/scaff/scaffui"

type StackProps struct {
	children []scaffui.NodeBuilder
}

func (sp *StackProps) Child(builder scaffui.NodeBuilder) {
	sp.children = append(sp.children, builder)
}

func Stack(create func(t *scaffui.Tracker, props *StackProps)) scaffui.NodeBuilder {

	// The default behavior of multi-node covers all of the things we want the stack to do, so no work for us :D
	return scaffui.UseMultiNode("stack", create, func(core *scaffui.MultiChildConstruct[StackProps]) {
		for _, child := range core.Props().children {
			core.Child(child)
		}
	})
}
