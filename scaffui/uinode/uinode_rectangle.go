package uinode

import (
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/paint"

	"image/color"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uispec"
	"github.com/Liphium/scaff/scath"
)

type RectangleProps struct {
	WantedConstraints scath.Constraints
	Padding           scath.Padding
	FillColor         color.Color
	BorderRadius      float64
	StrokeThickness   float64
	StrokeColor       color.Color

	*scaffui.AcceptChild
}

func Rectangle(create func(t *scaff.Tracker, props *RectangleProps)) scaffui.NodeBuilder {
	return scaffui.SingleNode(scaffui.SingleNodeCreate[RectangleProps]{
		ID: "rectangle",
		DefaultProps: RectangleProps{
			WantedConstraints: scath.Unconstrained(),
			Padding:           scath.Pad(0),
			FillColor:         color.White,
			BorderRadius:      0,
			StrokeColor:       color.Black,
			StrokeThickness:   0,
			AcceptChild:       &scaffui.AcceptChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.SingleChildProps[RectangleProps]) {
			props.OnWantedConstraints = func(node *scaffui.SingleChildNode[RectangleProps], _ scath.Constraints) scath.Constraints {
				return node.Props().WantedConstraints
			}

			props.OnLayout = func(node *scaffui.SingleChildNode[RectangleProps]) (scath.Vec, error) {
				props := node.Props()
				spec := uispec.SingleChildBoxSpec{
					Parent:  node.Constraints(),
					Wanted:  optional.With(props.WantedConstraints),
					Padding: props.Padding,
				}

				child, ok := node.Child()
				if ok {
					return spec.LayoutWithChild(child.Current())
				}

				return spec.LayoutWithoutChild()
			}

			props.OnDraw = func(node *scaffui.SingleChildNode[RectangleProps], position scath.Vec, painter paint.Painter) {
				props := node.Props()
				painter.Paint(paint.Rectangle{
					Position:     position.AddC(props.StrokeThickness),
					Size:         node.Size().AddC(-props.StrokeThickness * 2),
					FillColor:    props.FillColor,
					BorderRadius: props.BorderRadius,
				})
				if _, _, _, a := props.StrokeColor.RGBA(); a != 0 && props.StrokeThickness != 0 {
					painter.Paint(paint.RectangleStroke{
						Position:     position,
						Size:         node.Size(),
						Color:        node.Props().StrokeColor,
						BorderRadius: node.Props().BorderRadius,
						Thickness:    props.StrokeThickness,
					})
				}

				node.DrawChild(position.Add(props.Padding.ToVecTopLeft()), painter)
			}
		},
	})
}
