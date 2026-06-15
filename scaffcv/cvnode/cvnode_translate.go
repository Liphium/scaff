package cvnode

import (
	"time"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/engine"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
)

type TranslateProps struct {
	Offset func(now time.Time) scath.Vec
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
				return node.Props().Offset(time.Now()).Add(node.Children()[0].Position())
			}

			methods.OnDraw = func(node *scaffcv.StandardNode[TranslateProps], c *scaff.Context, painter engine.Painter) {
				before := painter.Transform()

				transform := painter.Transform()
				newPos := node.Props().Offset(c.Now())
				transform.CamX -= newPos.X
				transform.CamY -= newPos.Y

				painter.SetTransform(transform)
				node.DrawChildren(c, painter)
				painter.SetTransform(before)
			}
		},
	})
}
