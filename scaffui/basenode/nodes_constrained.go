package basenode

import (
	"github.com/Liphium/scaff/paint"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uispec"
	"github.com/Liphium/scaff/scath"
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
	return scaffui.CreateSingleNode("constrained", create, func(core *scaffui.SingleChildProps[ConstrainedProps]) {

		// Pass the child to the core node
		if child, ok := core.Props().child.Value(); ok {
			core.Child(child)
		}

		core.Layout(func(node *scaffui.SingleChildNode[ConstrainedProps]) (scath.Vec, error) {
			spec := uispec.SingleChildBoxSpec{
				Parent:        node.Constraints(),
				Wanted:        node.Props().constraints,
				scath.Padding: scaffui.Pad(0),
			}

			if child, ok := node.Child(); ok {
				return spec.LayoutWithChild(child.Current())
			}
			return spec.LayoutWithoutChild()
		})

		core.Draw(func(node *scaffui.SingleChildNode[ConstrainedProps], position scath.Vec, renderer paint.Painter) {
			node.DrawChild(position, renderer)
		})
	})
}
