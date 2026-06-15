package cvnode

import (
	"image/color"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/engine"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
)

type RectangleProps struct {
	Position        scath.Vec
	Size            scath.Vec
	FillColor       color.Color
	StrokeThickness float64
	StrokeColor     color.Color
	BorderRadius    float64

	*scaffcv.AcceptChild
}

func Rectangle(create func(t *scaff.Tracker, props *RectangleProps)) scaffcv.NodeBuilder {
	return scaffcv.Standard(scaffcv.StandardCreate[RectangleProps]{
		ID: "rectangle",
		DefaultProps: RectangleProps{
			Position:        scath.Zero,
			Size:            scath.Zero,
			FillColor:       color.White,
			BorderRadius:    0,
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

			methods.OnDraw = func(node *scaffcv.StandardNode[RectangleProps], c *scaff.Context, painter engine.Painter) {
				painter.Paint(engine.Rectangle{
					Position:     node.Props().Position,
					Size:         node.Props().Size,
					FillColor:    node.Props().FillColor,
					BorderRadius: node.Props().BorderRadius,
				})

				// Draw stroke when there
				if node.Props().StrokeThickness != 0 {
					painter.Paint(engine.RectangleStroke{
						Position:     node.Props().Position,
						Size:         node.Props().Size,
						Color:        node.Props().StrokeColor,
						BorderRadius: node.Props().BorderRadius,
					})
				}
			}
		},
	})
}
