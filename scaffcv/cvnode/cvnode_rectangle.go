package cvnode

import (
	"image/color"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
)

type RectangleProps struct {
	Position        scath.Vec
	Size            scath.Vec
	FillColor       color.Color
	StrokeThickness float64
	StrokeColor     color.Color
	borderRadius    float64

	*scaffcv.AcceptChild
}

func Rectangle(create func(t *scaff.Tracker, props *RectangleProps)) scaffcv.NodeBuilder {
	return scaffcv.Standard(scaffcv.StandardCreate[RectangleProps]{
		ID: "rectangle",
		DefaultProps: RectangleProps{
			Position:        scath.Zero,
			Size:            scath.Vec{X: 100, Y: 100},
			FillColor:       color.White,
			borderRadius:    0,
			StrokeColor:     color.Black,
			StrokeThickness: 0,
			AcceptChild:     &scaffcv.AcceptChild{},
		},
		PropsCreator: create,
		Create: func(methods *scaffcv.StandardMethods[RectangleProps]) {
			methods.Position = func(node *scaffcv.StandardNode[RectangleProps]) scath.Vec {
				return node.Props().Position
			}
			methods.Size = func(node *scaffcv.StandardNode[RectangleProps]) scath.Vec {
				return node.Props().Size
			}

			methods.OnDraw = func(node *scaffcv.StandardNode[RectangleProps], c *scaff.Context, painter paint.Painter) {
				painter.Paint(paint.Rectangle{
					Position:     node.Props().Position,
					Size:         node.Props().Size,
					FillColor:    node.Props().FillColor,
					BorderRadius: node.Props().borderRadius,
				})

				// Draw stroke when there
				if node.Props().StrokeThickness != 0 {
					painter.Paint(paint.RectangleStroke{
						Position:     node.Props().Position,
						Size:         node.Props().Size,
						Color:        node.Props().StrokeColor,
						BorderRadius: node.Props().borderRadius,
					})
				}
			}
		},
	})
}
