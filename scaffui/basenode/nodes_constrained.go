package basenode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uispec"
	"github.com/Liphium/scaff/smath"
)

type ConstrainedProps struct {
	child       optional.O[scaffui.NodeBuilder]
	constraints optional.O[scaffui.Constraints]
}

func (cp *ConstrainedProps) Constraints(constraints scaffui.Constraints) {
	cp.constraints.SetValue(constraints)
}

func (cp *ConstrainedProps) Child(builder scaffui.NodeBuilder) {
	cp.child.SetValue(builder)
}

func Constrained(create func(t *scaff.Tracker, props *ConstrainedProps)) scaffui.NodeBuilder {
	return scaffui.UseSingleNode("constrained", create, func(core *scaffui.SingleChildConstruct[ConstrainedProps]) {

		// Pass the child to the core node
		if child, ok := core.Props().child.Value(); ok {
			core.Child(child)
		}

		core.Layout(func(node *scaffui.SingleChildNode[ConstrainedProps]) (scaffui.Size, error) {
			spec := uispec.SingleChildBoxSpec{
				Parent:  node.Constraints(),
				Wanted:  node.Props().constraints,
				Padding: scaffui.Pad(0),
			}

			if child, ok := node.Child(); ok {
				return spec.LayoutWithChild(child.Current())
			}
			return spec.LayoutWithoutChild()
		})

		core.Draw(func(node *scaffui.SingleChildNode[ConstrainedProps], position smath.Vec, renderer scaffui.Renderer) {
			node.DrawChild(position, renderer)
		})
	})
}
