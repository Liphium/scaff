package basenode

import (
	"image/color"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scath"
)

type TextProps struct {
	text     string
	font     string
	fontSize float64
	color    color.Color
	*scaffui.AcceptNoChild
}

func (tp *TextProps) Text(text string) {
	tp.text = text
}

func (tp *TextProps) Font(font string) {
	tp.font = font
}

func (tp *TextProps) FontSize(size float64) {
	tp.fontSize = size
}

func (tp *TextProps) Color(color color.RGBA) {
	tp.color = color
}

func Text(create func(t *scaff.Tracker, props *TextProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode(scaffui.SingleNodeCreate[TextProps]{
		ID: "text",
		DefaultProps: TextProps{
			text:          "Scaff",
			fontSize:      16,
			color:         color.White,
			AcceptNoChild: &scaffui.AcceptNoChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.SingleChildProps[TextProps]) {
			props.WantedConstraints(func(node *scaffui.SingleChildNode[TextProps], parent scath.Constraints) scath.Constraints {
				return scath.Unconstrained()
			})

			props.Draw(func(node *scaffui.SingleChildNode[TextProps], position scath.Vec, painter paint.Painter) {
				painter.Paint(paint.Text{
					Font:     node.Props().font,
					Text:     node.Props().text,
					Color:    node.Props().color,
					FontSize: node.Props().fontSize,
				})
			})
		},
	})
}
