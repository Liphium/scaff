package scaffui

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type ViewportProps struct {
	assets *paint.AssetManager
	child  optional.O[NodeBuilder]
}

func (vp *ViewportProps) Child(builder NodeBuilder) {
	vp.child.SetValue(builder)
}

func (vp ViewportProps) GetBuilders() []scaff.NodeBuilder {
	return nil
}

// Viewport creates a viewport node that can be used to essentially mount a
func Viewport(assets *paint.AssetManager, create func(t *scaff.Tracker, props *ViewportProps)) scaff.NodeBuilder {
	return scaff.CreateSingleNode("viewport", create, func(props *scaff.SingleChildProps[ViewportProps]) {
		var root *MountedNode
		var renderer *paint.EbitenPainter

		props.Load(func(node *scaff.SingleChildNode[ViewportProps], parent scaff.Node) {
			child, ok := node.Props().child.Value()
			if !ok {
				return
			}

			root = NewMountedFromBuilder(child)
			root.Load(nil)
		})

		props.Unload(func(node *scaff.SingleChildNode[ViewportProps]) {
			if root != nil {
				root.Unload()
			}
		})

		props.Draw(func(node *scaff.SingleChildNode[ViewportProps], c *scaff.Context, screen *ebiten.Image) {
			if root == nil {
				return
			}

			// Initialize the UI and stuff
			if renderer == nil {
				renderer = paint.NewEbitenPainter(ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy()), true, assets)
				root.Current().SetConstraints(scath.Loose(float64(c.Width), float64(c.Height)))
				_, err := root.Current().Layout()
				if err != nil {
					log.Error("layout error", "err", err)
				}
			}

			// Update all of the stuff
			result, err := root.Update(nil, c)
			if result.SizeChanged || err != nil {
				log.Warn("relayout or error happend", "result", result, "err", err)
				return
			}

			// Draw the stuff
			renderer.Clear()
			root.Current().Draw(scath.Zero, renderer)

			screen.DrawImage(renderer.Screen(), &ebiten.DrawImageOptions{})
		})

	})
}
