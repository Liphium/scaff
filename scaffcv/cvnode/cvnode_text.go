package cvnode

import (
	"image/color"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TextProps struct {
	Text           string
	Position       scath.Vec
	Color          color.Color
	PrimaryAlign   text.Align
	SecondaryAlign text.Align
	scaffcv.AcceptNoChild
}

// Text creates a simple Text node with a Position and Text.
func Text(create func(t *scaff.Tracker, props *TextProps)) scaffcv.NodeBuilder {
	return scaffcv.SingleNode(scaffcv.SingleNodeCreate[TextProps]{
		ID: "Text",
		DefaultProps: TextProps{
			Text:           "",
			Position:       scath.Vec{},
			Color:          color.White,
			PrimaryAlign:   text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
		},
		PropsCreator: create,
		Create: func(props *scaffcv.SingleChildProps[TextProps]) {
			props.OnDraw = func(node *scaffcv.SingleChildNode[TextProps], c *scaff.Context, painter paint.Painter) {
				painter.Paint(paint.Text{
					Text:           node.Props().Text,
					FontSize:       20,
					Position:       node.Props().Position,
					Color:          node.Props().Color,
					PrimaryAlign:   node.Props().PrimaryAlign,
					SecondaryAlign: node.Props().SecondaryAlign,
				})
			}
		},
	})
}
