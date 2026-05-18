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
	constraints optional.O[scath.Constraints]
	*scaffui.AcceptChild
}

func (cp *ConstrainedProps) Constraints(constraints scath.Constraints) {
	cp.constraints.SetValue(constraints)
}

func Constrained(create func(t *scaff.Tracker, props *ConstrainedProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode("constrained", create, func(core *scaffui.SingleChildProps[ConstrainedProps]) {

		core.Layout(func(node *scaffui.SingleChildNode[ConstrainedProps]) (scath.Vec, error) {
			spec := uispec.SingleChildBoxSpec{
				Parent:  node.Constraints(),
				Wanted:  node.Props().constraints,
				Padding: scath.Pad(0),
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
