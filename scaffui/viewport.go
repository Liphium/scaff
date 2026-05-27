package scaffui

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type ViewportProps struct {
	child optional.O[NodeBuilder]
}

func (vp *ViewportProps) Child(builder NodeBuilder) {
	vp.child.SetValue(builder)
}

func (vp ViewportProps) GetBuilders() []scaff.NodeBuilder {
	return nil
}

// Viewport creates a viewport node that can be used to essentially mount a
func Viewport(create func(t *scaff.Tracker, props *ViewportProps)) scaff.NodeBuilder {
	return scaff.SingleNode("viewport", create, func(props *scaff.SingleChildProps[ViewportProps]) {
		var root *MountedNode
		var renderer *paint.EbitenPainter

		props.Load(func(node *scaff.SingleChildNode[ViewportProps], parent scaff.Node) {
			child, ok := node.Props().child.Value()
			if !ok {
				return
			}

			root = NewMountedFromBuilder(child, node.Context())
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
			firstRender := false
			if renderer == nil {
				renderer = paint.NewEbitenPainter(ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy()), true, node.Context().AssetManager())
				root.Current().SetConstraints(scath.Loose(float64(c.Width), float64(c.Height)))
				_, err := root.Current().Layout()
				if err != nil {
					log.Error("layout error", "err", err)
				}
				firstRender = true
			}

			// Update all of the stuff
			result, err := root.Update(nil, c, node.Context())
			if result.SizeChanged || err != nil {
				log.Warn("relayout or error happend", "result", result, "err", err)
				return
			}

			// Draw the stuff (only if changed)
			if result.AnythingChanged || firstRender {
				renderer.Clear()
				root.Current().Draw(scath.Zero, renderer)
			}

			screen.DrawImage(renderer.Screen(), &ebiten.DrawImageOptions{})
		})

		props.HandleEvent(func(node *scaff.SingleChildNode[ViewportProps], c *scaff.Context, event scaff.Event) error {

			// When the size changes, delete the renderer (screen size changes and stuff and we can't use the old stuff anymore anyway)
			if event.EventID() == scaff.EventIdSizeChange {
				renderer = nil
			}

			if root != nil {
				root.Current().HandleEvent(c, event)
			}
			return nil
		})
	})
}
