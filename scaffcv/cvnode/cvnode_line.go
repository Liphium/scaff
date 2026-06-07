package cvnode

import (
	"image/color"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
)

type LineProps struct {
	Start     scath.Vec
	End       scath.Vec
	Color     color.Color
	Thickness float64

	scaffcv.AcceptNoChild
}

func Line(create func(t *scaff.Tracker, props *LineProps)) scaffcv.NodeBuilder {
	return scaffcv.Standard(scaffcv.StandardCreate[LineProps]{
		ID: "line",
		DefaultProps: LineProps{
			Thickness: 4,
			Color:     color.White,
		},
		PropsCreator: create,
		Create: func(methods *scaffcv.StandardMethods[LineProps]) {
			methods.OnDraw = func(node *scaffcv.StandardNode[LineProps], c *scaff.Context, painter paint.Painter) {
				painter.Paint(paint.Line{
					Start:     node.Props().Start,
					End:       node.Props().End,
					Color:     node.Props().Color,
					Thickness: node.Props().Thickness,
				})
			}
		},
	})
}
