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
	wantedConstraints scath.Constraints
	padding           scath.Padding
	fillColor         color.Color
	borderRadius      float64
	strokeThickness   float64
	strokeColor       color.Color

	*scaffui.AcceptChild
}

func (rp *RectangleProps) WantedConstraints(constraints scath.Constraints) {
	rp.wantedConstraints = constraints
}

func (rp *RectangleProps) Padding(padding scath.Padding) {
	rp.padding = padding
}

func (rp *RectangleProps) FillColor(fillColor color.Color) {
	rp.fillColor = fillColor
}

func (rp *RectangleProps) BorderRadius(borderRadius float64) {
	rp.borderRadius = borderRadius
}

func (r *RectangleProps) StrokeColor(color color.Color) {
	r.strokeColor = color
}

func (r *RectangleProps) StrokeThickness(thickness float64) {
	r.strokeThickness = thickness
}

func Rectangle(create func(t *scaff.Tracker, props *RectangleProps)) scaffui.NodeBuilder {
	return scaffui.SingleNode(scaffui.SingleNodeCreate[RectangleProps]{
		ID: "rectangle",
		DefaultProps: RectangleProps{
			wantedConstraints: scath.Unconstrained(),
			padding:           scath.Pad(0),
			fillColor:         color.White,
			borderRadius:      0,
			strokeColor:       color.Black,
			strokeThickness:   0,
			AcceptChild:       &scaffui.AcceptChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.SingleChildProps[RectangleProps]) {
			props.WantedConstraints(func(node *scaffui.SingleChildNode[RectangleProps], _ scath.Constraints) scath.Constraints {
				return node.Props().wantedConstraints
			})

			props.Layout(func(node *scaffui.SingleChildNode[RectangleProps]) (scath.Vec, error) {
				props := node.Props()
				spec := uispec.SingleChildBoxSpec{
					Parent:  node.Constraints(),
					Wanted:  optional.With(props.wantedConstraints),
					Padding: props.padding,
				}

				child, ok := node.Child()
				if ok {
					return spec.LayoutWithChild(child.Current())
				}

				return spec.LayoutWithoutChild()
			})

			props.Draw(func(node *scaffui.SingleChildNode[RectangleProps], position scath.Vec, painter paint.Painter) {
				props := node.Props()
				painter.Paint(paint.Rectangle{
					Position:     position.AddC(props.strokeThickness),
					Size:         node.Size().AddC(-props.strokeThickness * 2),
					FillColor:    props.fillColor,
					BorderRadius: props.borderRadius,
				})
				if _, _, _, a := props.strokeColor.RGBA(); a != 0 && props.strokeThickness != 0 {
					painter.Paint(paint.RectangleStroke{
						Position:     position,
						Size:         node.Size(),
						Color:        node.Props().strokeColor,
						BorderRadius: node.Props().borderRadius,
						Thickness:    props.strokeThickness,
					})
				}

				node.DrawChild(position.Add(props.padding.ToVecTopLeft()), painter)
			})
		},
	})
}
