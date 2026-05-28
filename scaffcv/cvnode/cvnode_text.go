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
	text           string
	position       scath.Vec
	color          color.Color
	primaryAlign   text.Align
	secondaryAlign text.Align
	*scaffcv.AcceptNoChild
}

func (tp *TextProps) Text(text string) {
	tp.text = text
}

func (tp *TextProps) Position(pos scath.Vec) {
	tp.position = pos
}

func (tp *TextProps) Color(c color.Color) {
	tp.color = c
}

func (tp *TextProps) PrimaryAlign(align text.Align) {
	tp.primaryAlign = align
}

func (tp *TextProps) SecondaryAlign(align text.Align) {
	tp.secondaryAlign = align
}

// Text creates a simple text node with a position and text.
func Text(create func(t *scaff.Tracker, props *TextProps)) scaffcv.NodeBuilder {
	return scaffcv.SingleNode(scaffcv.SingleNodeCreate[TextProps]{
		ID: "text",
		DefaultProps: TextProps{
			text:           "",
			position:       scath.Vec{},
			color:          color.White,
			primaryAlign:   text.AlignCenter,
			secondaryAlign: text.AlignCenter,
			AcceptNoChild:  &scaffcv.AcceptNoChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffcv.SingleChildProps[TextProps]) {
			props.Draw(func(node *scaffcv.SingleChildNode[TextProps], c *scaff.Context, painter paint.Painter) {
				painter.Paint(paint.Text{
					Text:           node.Props().text,
					FontSize:       20,
					Position:       node.Props().position,
					Color:          node.Props().color,
					PrimaryAlign:   node.Props().primaryAlign,
					SecondaryAlign: node.Props().secondaryAlign,
				})
			})
		},
	})
}

