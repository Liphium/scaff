package cvnode

import (
	"image/color"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
)

type RectangleProps struct {
	position        scath.Vec
	size            scath.Vec
	fillColor       color.Color
	strokeThickness float64
	strokeColor     color.Color
	borderRadius    float64

	*scaffcv.AcceptChild
}

func (r *RectangleProps) Position(position scath.Vec) {
	r.position = position
}

func (r *RectangleProps) Size(size scath.Vec) {
	r.size = size
}

func (r *RectangleProps) FillColor(color color.Color) {
	r.fillColor = color
}

func (r *RectangleProps) StrokeColor(color color.Color) {
	r.strokeColor = color
}

func (r *RectangleProps) StrokeThickness(thickness float64) {
	r.strokeThickness = thickness
}

func Rectangle(create func(t *scaff.Tracker, props *RectangleProps)) scaffcv.NodeBuilder {
	return scaffcv.SingleNode(scaffcv.SingleNodeCreate[RectangleProps]{
		ID: "rectangle",
		DefaultProps: RectangleProps{
			position:        scath.Zero,
			size:            scath.Vec{X: 100, Y: 100},
			fillColor:       color.White,
			borderRadius:    0,
			strokeColor:     color.Black,
			strokeThickness: 0,
			AcceptChild:     &scaffcv.AcceptChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffcv.SingleChildProps[RectangleProps]) {
			props.Draw(func(node *scaffcv.SingleChildNode[RectangleProps], c *scaff.Context, painter paint.Painter) {
				painter.Paint(paint.Rectangle{
					Position:     node.Props().position,
					Size:         node.Props().size,
					FillColor:    node.Props().fillColor,
					BorderRadius: node.Props().borderRadius,
				})

				// Draw stroke when there
				if node.Props().strokeThickness != 0 {
					painter.Paint(paint.RectangleStroke{
						Position:     node.Props().position,
						Size:         node.Props().size,
						Color:        node.Props().strokeColor,
						BorderRadius: node.Props().borderRadius,
					})
				}
			})
		},
	})
}
