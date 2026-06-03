package scaffui

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type ViewportProps struct {
	child *AcceptChild
	*scaff.AcceptNoChild
}

func (vp ViewportProps) Child(builder NodeBuilder) {
	vp.child.Child(builder)
}

// Viewport creates a viewport node that can be used to essentially mount a
func Viewport(create func(t *scaff.Tracker, props *ViewportProps)) scaff.NodeBuilder {

	// TODO: For viewport:
	// - Add checks to the update queue if there are updates (then use that to determine if we even need to call Update() in the first place)

	return scaff.Standard(scaff.StandardCreate[ViewportProps]{
		ID: "viewport",
		DefaultProps: ViewportProps{
			child:         &AcceptChild{},
			AcceptNoChild: &scaff.AcceptNoChild{},
		},
		PropsCreator: create,
		Create: func(props *scaff.StandardMethods[ViewportProps]) {
			var root Node
			var context *scaff.BuildContext
			var renderer *paint.EbitenPainter

			props.OnLoad = func(node *scaff.StandardNode[ViewportProps], parent scaff.Node) {
				context = node.Context().CopyWithNewUpdateQueue()
			}

			props.OnPropsChanged = func(node *scaff.StandardNode[ViewportProps]) {
				changed := node.Props().child.GetChanged()
				if len(changed) > 0 {
					root = node.Props().child.GetBuilders()[0](context)
					root.Load(nil)
					node.Props().child.ClearChanged()
				}
			}

			props.OnUnload = func(node *scaff.StandardNode[ViewportProps]) {
				if root != nil {
					root.Unload()
				}
			}

			props.OnDraw = func(node *scaff.StandardNode[ViewportProps], c *scaff.Context, screen *ebiten.Image) {
				if root == nil {
					return
				}

				// Initialize the UI and stuff
				firstRender := false
				if renderer == nil {
					renderer = paint.NewEbitenPainter(ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy()), true, node.Context().AssetManager())
					root.SetConstraints(scath.Loose(float64(c.Width()), float64(c.Height())))
					_, err := root.Layout()
					if err != nil {
						log.Error("layout error", "err", err)
					}
					firstRender = true
				}

				// Execute update queue
				context.UpdateQueue().Update()

				// Update all of the stuff
				result, err := root.Update()
				if result.SizeChanged || err != nil {
					log.Warn("relayout or error happend", "result", result, "err", err)
					return
				}

				// Draw the stuff (only if changed)
				if result.AnythingChanged || firstRender {
					renderer.Clear()
					root.Draw(scath.Zero, renderer)
				}

				screen.DrawImage(renderer.Screen(), &ebiten.DrawImageOptions{})
			}

			props.OnHandleEvent = func(node *scaff.StandardNode[ViewportProps], c *scaff.Context, event scaff.Event) error {

				// When the size changes, delete the renderer (screen size changes and stuff and we can't use the old stuff anymore anyway)
				if event.EventID() == scaff.EventIdSizeChange {
					renderer = nil
				}

				if root != nil {
					root.HandleEvent(c, event)
				}
				return nil
			}
		},
	})
}
