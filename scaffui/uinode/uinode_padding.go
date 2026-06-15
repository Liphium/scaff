package uinode

import (
	"github.com/Liphium/scaff/engine"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uispec"
	"github.com/Liphium/scaff/scath"
)

var _ scaff.ChildProps[scaffui.NodeBuilder] = PaddingProps{}

type PaddingProps struct {
	Padding optional.O[scath.Padding]
	*scaffui.AcceptChild
}

func Padding(create func(t *scaff.Tracker, props *PaddingProps)) scaffui.NodeBuilder {
	return scaffui.Standard(scaffui.StandardCreate[PaddingProps]{
		ID: "Padding",
		DefaultProps: PaddingProps{
			AcceptChild: &scaffui.AcceptChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.StandardMethods[PaddingProps]) {
			// In Layout, make sure to give the child less constraints (subtracted by Padding, handled by uispec)
			props.OnLayout = func(node *scaffui.StandardNode[PaddingProps]) (scath.Vec, error) {
				spec := uispec.SingleChildBoxSpec{
					Parent:  node.Constraints(),
					Wanted:  optional.None[scath.Constraints](),
					Padding: node.Props().Padding.Or(scath.Pad(0)),
				}

				if child, ok := node.Child(); ok {
					return spec.LayoutWithChild(child)
				}
				return spec.LayoutWithoutChild()
			}

			// Draw child at padded position
			props.OnDraw = func(node *scaffui.StandardNode[PaddingProps], position scath.Vec, renderer engine.Painter) {
				if child, ok := node.Child(); ok {
					child.Draw(position.Add(node.Props().Padding.Or(scath.Pad(0)).ToVecTopLeft()), renderer)
				}
			}
		},
	})
}
