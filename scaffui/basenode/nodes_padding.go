package basenode

import (
	"github.com/Liphium/scaff/paint"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uispec"
	"github.com/Liphium/scaff/scath"
)

var _ scaff.ChildProps[scaffui.NodeBuilder] = PaddingProps{}

type PaddingProps struct {
	padding optional.O[scath.Padding]
	*scaffui.AcceptChild
}

func (pp *PaddingProps) Padding(padding scath.Padding) {
	pp.padding.SetValue(padding)
}

func Padding(create func(t *scaff.Tracker, props *PaddingProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode("padding", create, func(props *scaffui.SingleChildProps[PaddingProps]) {

		// In Layout, make sure to give the child less constraints (subtracted by padding, handled by uispec)
		props.Layout(func(node *scaffui.SingleChildNode[PaddingProps]) (scath.Vec, error) {
			spec := uispec.SingleChildBoxSpec{
				Parent:  node.Constraints(),
				Wanted:  optional.None[scath.Constraints](),
				Padding: node.Props().padding.Or(scath.Pad(0)),
			}

			if child, ok := node.Child(); ok {
				return spec.LayoutWithChild(child.Current())
			}
			return spec.LayoutWithoutChild()
		})

		// Draw child at padded position
		props.Draw(func(node *scaffui.SingleChildNode[PaddingProps], position scath.Vec, renderer paint.Painter) {
			node.DrawChild(position.Add(node.Props().padding.Or(scath.Pad(0)).ToVecTopLeft()), renderer)
		})
	})
}
