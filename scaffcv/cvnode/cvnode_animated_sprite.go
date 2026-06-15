package cvnode

import (
	"image"
	"time"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/engine"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type AnimatedSpriteProps struct {
	Position scath.Vec

	SpriteSheet string
	FrameWidth  int
	FrameHeight int

	FramesToPlay []TilePosition
	Duration     time.Duration

	scaffcv.AcceptNoChild
}

func AnimatedSprite(create func(t *scaff.Tracker, props *AnimatedSpriteProps)) scaffcv.NodeBuilder {
	return scaffcv.Standard(scaffcv.StandardCreate[AnimatedSpriteProps]{
		ID: "animated-sprite",
		DefaultProps: AnimatedSpriteProps{
			Position:     scath.Zero,
			FrameWidth:   16,
			FrameHeight:  16,
			FramesToPlay: []TilePosition{},
			Duration:     time.Millisecond * 500,
		},
		PropsCreator: create,
		Create: func(methods *scaffcv.StandardMethods[AnimatedSpriteProps]) {
			methods.Position = func(node *scaffcv.StandardNode[AnimatedSpriteProps]) scath.Vec {
				return node.Props().Position
			}
			methods.Size = func(node *scaffcv.StandardNode[AnimatedSpriteProps]) scath.Vec {
				return scath.Vec{X: float64(node.Props().FrameWidth), Y: float64(node.Props().FrameHeight)}
			}

			// Properties for the current information
			start := time.Now()
			totalDuration := time.Millisecond // Will get set in PropsChanged

			methods.OnPropsChanged = func(node *scaffcv.StandardNode[AnimatedSpriteProps]) {
				totalDuration = time.Duration(len(node.Props().FramesToPlay)) * node.Props().Duration
			}

			methods.OnDraw = func(node *scaffcv.StandardNode[AnimatedSpriteProps], c *scaff.Context, painter engine.Painter) {
				spritesheet, err := node.Context().AssetManager().GetImage(node.Props().SpriteSheet)
				if err != nil {
					log.Error("couldn't load image for spritesheet", "i", node.Props().SpriteSheet)
					return
				}

				currentFrame := (c.Now().Sub(start) % totalDuration) / node.Props().Duration
				frame := node.Props().FramesToPlay[currentFrame]

				frameImage := spritesheet.SubImage(image.Rectangle{
					Min: image.Pt(frame.X, frame.Y),
					Max: image.Pt(frame.X+node.Props().FrameWidth, frame.Y+node.Props().FrameHeight),
				})

				op := &ebiten.DrawImageOptions{
					Filter: ebiten.FilterPixelated,
				}
				op.GeoM.Translate(node.Props().Position.X, node.Props().Position.Y)

				painter.DrawRaw(frameImage.(*ebiten.Image), op)
			}
		},
	})
}
