package scaffcv

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type CanvasProps struct {
	// The current camera position
	Position scath.Vec

	// The current zoom factor
	Zoom float64

	// The current smooth type
	SmoothType SmoothType

	// The current camera smoothing options
	SmoothOptions *SmoothOptions

	child *AcceptChild
	*scaff.AcceptNoChild
}

func (vp CanvasProps) Child(builder NodeBuilder) {
	vp.child.Child(builder)
}

func Canvas(create func(t *scaff.Tracker, props *CanvasProps)) scaff.NodeBuilder {
	return scaff.Standard(scaff.StandardCreate[CanvasProps]{
		ID: "canvas",
		DefaultProps: CanvasProps{
			SmoothType:    CameraMovementNone,
			SmoothOptions: DefaultSmoothOptions(),
			child:         &AcceptChild{},
			AcceptNoChild: &scaff.AcceptNoChild{},
		},
		PropsCreator: create,
		Create: func(props *scaff.StandardMethods[CanvasProps]) {
			var context *BuildContext
			var cv CanvasProps
			var root Node
			var painter *paint.EbitenPainter
			sizeUpdate := true

			props.OnLoad = func(node *scaff.StandardNode[CanvasProps], parent scaff.Node) {
				context = &BuildContext{
					BuildContext: node.Context(),
				}
			}

			props.OnPropsChanged = func(node *scaff.StandardNode[CanvasProps]) {
				changed := node.Props().child.GetChanged()
				if len(changed) > 0 {
					root = node.Props().child.GetBuilders()[0](context)
					root.Load(nil)
					node.Props().child.ClearChanged()
				}

				if context.cam != nil {
					context.cam.SmoothType = node.Props().SmoothType
					context.cam.SmoothOptions = node.Props().SmoothOptions
				}
			}

			props.OnUpdate = func(node *scaff.StandardNode[CanvasProps], c *scaff.Context) error {
				if root == nil {
					return nil
				}

				// Forward updates to the child nodes
				return root.Update(c)
			}

			props.OnDraw = func(node *scaff.StandardNode[CanvasProps], c *scaff.Context, image *ebiten.Image) {
				if root == nil {
					return
				}

				// This is run on the first frame as well to create the painter
				if sizeUpdate {
					screen := ebiten.NewImage(image.Bounds().Dx(), image.Bounds().Dy())
					painter = paint.NewEbitenPainter(screen, true, node.Context().AssetManager())

					// Create new camera (old one will only have invalid positions)
					context.cam = NewCamera(cv.Position.X, cv.Position.Y, c.Width(), c.Height())
					context.cam.SmoothType = cv.SmoothType
					context.cam.SmoothOptions = cv.SmoothOptions
					sizeUpdate = false
				}

				// Set the proper position on the camera (this needs to be called every update so the smoothing is updated)
				context.cam.LookAt(node.Props().Position.X, node.Props().Position.Y)
				context.cam.ZoomFactor = node.Props().Zoom

				// Actually draw the root
				painter.Clear()
				painter.SetTransform(paint.Transform{
					CamX:          context.cam.X,
					CamY:          context.cam.Y,
					CenterOffsetX: context.cam.CenterOffsetX,
					CenterOffsetY: context.cam.CenterOffsetY,
					Angle:         context.cam.Angle,
					ZoomFactor:    context.cam.ZoomFactor,
				})
				root.Draw(c, painter)
				image.DrawImage(painter.Screen(), &ebiten.DrawImageOptions{})
			}

			props.OnHandleEvent = func(node *scaff.StandardNode[CanvasProps], c *scaff.Context, event scaff.Event) error {
				if root == nil {
					return nil
				}

				// For a size change, make sure we properly handle it
				if event.EventID() == scaff.EventIdSizeChange {
					sizeUpdate = true
				}

				// Forward events to the root
				return root.HandleEvent(c, event)
			}
		},
	})
}
