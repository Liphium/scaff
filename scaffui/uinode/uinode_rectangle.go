package uinode

import (
	"github.com/Liphium/scaff/paint"

	"image/color"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uispec"
	"github.com/Liphium/scaff/scath"
)

type RectangleProps struct {
	wantedConstraints optional.O[scath.Constraints]
	padding           optional.O[scath.Padding]
	fillColor         optional.O[color.RGBA]
	borderRadius      optional.O[int]
	*scaffui.AcceptChild
}

func (rp *RectangleProps) WantedConstraints(constraints scath.Constraints) {
	rp.wantedConstraints.SetValue(constraints)
}

func (rp *RectangleProps) Padding(padding scath.Padding) {
	rp.padding.SetValue(padding)
}

func (rp *RectangleProps) FillColor(fillColor color.RGBA) {
	rp.fillColor.SetValue(fillColor)
}

func (rp *RectangleProps) BorderRadius(borderRadius int) {
	rp.borderRadius.SetValue(borderRadius)
}

func Rectangle(create func(t *scaff.Tracker, props *RectangleProps)) scaffui.NodeBuilder {
	return scaffui.SingleNode(scaffui.SingleNodeCreate[RectangleProps]{
		ID: "rectangle",
		DefaultProps: RectangleProps{
			AcceptChild: &scaffui.AcceptChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.SingleChildProps[RectangleProps]) {
			props.WantedConstraints(func(node *scaffui.SingleChildNode[RectangleProps], _ scath.Constraints) scath.Constraints {
				return node.Props().wantedConstraints.Or(scath.Unconstrained())
			})

			props.Layout(func(node *scaffui.SingleChildNode[RectangleProps]) (scath.Vec, error) {
				props := node.Props()
				spec := uispec.SingleChildBoxSpec{
					Parent:  node.Constraints(),
					Wanted:  props.wantedConstraints,
					Padding: props.padding.Or(scath.Pad(0)),
				}

				child, ok := node.Child()
				if ok {
					return spec.LayoutWithChild(child.Current())
				}

				return spec.LayoutWithoutChild()
			})

			props.Draw(func(node *scaffui.SingleChildNode[RectangleProps], position scath.Vec, renderer paint.Painter) {
				props := node.Props()
				renderer.Paint(paint.Rectangle{
					Position:     position,
					Size:         node.Size(),
					FillColor:    props.fillColor.Or(color.RGBA{255, 255, 255, 255}),
					BorderRadius: props.borderRadius.Or(0),
				})

				node.DrawChild(position.Add(props.padding.Or(scath.Pad(0)).ToVecTopLeft()), renderer)
			})
		},
	})
}
