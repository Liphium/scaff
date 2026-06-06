package cvnode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
)

type TranslateProps struct {
	Offset scath.Vec
	*scaffcv.AcceptChild
}

func Translate(create func(t *scaff.Tracker, props *TranslateProps)) scaffcv.NodeBuilder {
	return scaffcv.Standard(scaffcv.StandardCreate[TranslateProps]{
		ID: "translate",
		DefaultProps: TranslateProps{
			AcceptChild: scaffcv.EmptyChild(),
		},
		PropsCreator: create,
		Create: func(methods *scaffcv.StandardMethods[TranslateProps]) {
			methods.Position = func(node *scaffcv.StandardNode[TranslateProps]) scath.Vec {
				if len(node.Children()) == 0 {
					return scath.Zero
				}
				return node.Props().Offset.Add(node.Children()[0].Position())
			}

			methods.OnDraw = func(node *scaffcv.StandardNode[TranslateProps], c *scaff.Context, painter paint.Painter) {
				before := painter.Transform()

				transform := painter.Transform()
				transform.CamX -= node.Props().Offset.X
				transform.CamY -= node.Props().Offset.Y

				painter.SetTransform(transform)
				node.DrawChildren(c, painter)
				painter.SetTransform(before)
			}
		},
	})
}
