package cvnode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
)

// Text creates a simple text node with a position and text.
func Text(text string, pos scath.Vec) scaffcv.NodeBuilder {
	return scaffcv.UseNode("text", func(props *scaffcv.SingleChildProps[any]) {
		props.Draw(func(node *scaffcv.SingleChildNode[any], c *scaff.Context, painter paint.Painter) {
			painter.Paint(paint.Text{
				Text:     text,
				Position: pos,
			})
		})
	})
}