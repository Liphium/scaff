package cvnode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type ImageProps struct {
	Path       string
	FilterMode ebiten.Filter
	Position   scath.Vec
	Size       scath.Vec
	scaffcv.AcceptNoChild
}

func Image(create func(t *scaff.Tracker, props *ImageProps)) scaffcv.NodeBuilder {
	return scaffcv.Standard(scaffcv.StandardCreate[ImageProps]{
		ID:           "image",
		DefaultProps: ImageProps{},
		PropsCreator: create,
		Create: func(methods *scaffcv.StandardMethods[ImageProps]) {
			methods.Position = func(node *scaffcv.StandardNode[ImageProps]) scath.Vec {
				return node.Props().Position
			}
			methods.Size = func(node *scaffcv.StandardNode[ImageProps]) scath.Vec {
				return node.Props().Size
			}

			methods.OnDraw = func(node *scaffcv.StandardNode[ImageProps], c *scaff.Context, painter paint.Painter) {
				painter.Paint(paint.Image{
					Position:   node.Props().Position,
					Size:       node.Props().Size,
					Path:       node.Props().Path,
					FilterMode: node.Props().FilterMode,
				})
			}
		},
	})
}
