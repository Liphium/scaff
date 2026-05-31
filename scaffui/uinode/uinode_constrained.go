package uinode

import (
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/paint"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uispec"
	"github.com/Liphium/scaff/scath"
)

type ConstrainedProps struct {
	Constraints optional.O[scath.Constraints]
	*scaffui.AcceptChild
}

func Constrained(create func(t *scaff.Tracker, props *ConstrainedProps)) scaffui.NodeBuilder {
	return scaffui.SingleNode(scaffui.SingleNodeCreate[ConstrainedProps]{
		ID: "constrained",
		DefaultProps: ConstrainedProps{
			Constraints: optional.None[scath.Constraints](),
			AcceptChild: &scaffui.AcceptChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.SingleChildProps[ConstrainedProps]) {
			props.OnLayout = func(node *scaffui.SingleChildNode[ConstrainedProps]) (scath.Vec, error) {
				spec := uispec.SingleChildBoxSpec{
					Parent:  node.Props().Constraints.Or(node.Constraints()),
					Wanted:  node.Props().Constraints,
					Padding: scath.Pad(0),
				}

				if child, ok := node.Child(); ok {
					return spec.LayoutWithChild(child.Current())
				}
				return spec.LayoutWithoutChild()
			}

			props.OnDraw = func(node *scaffui.SingleChildNode[ConstrainedProps], position scath.Vec, renderer paint.Painter) {
				node.DrawChild(position, renderer)
			}
		},
	})
}
