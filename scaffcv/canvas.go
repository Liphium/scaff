package scaffcv

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type CanvasProps struct {
	Position      scath.Vec
	SmoothType    SmoothType
	SmoothOptions SmoothOptions

	children []NodeBuilder
}

func (sp *CanvasProps) Child(builder NodeBuilder) {
	sp.children = append(sp.children, builder)
}

func (sp CanvasProps) GetBuilders() []NodeBuilder {
	return sp.children
}

func Canvas(create func(t *scaff.Tracker, props *CanvasProps)) scaff.NodeBuilder {
	return scaff.SingleNode(scaff.SingleNodeCreate[*scaff.AcceptNoChild]{
		ID:           "canvas",
		DefaultProps: &scaff.AcceptNoChild{},
		Create: func(props *scaff.SingleChildProps[*scaff.AcceptNoChild]) {
			var camera *Camera
			var cv CanvasProps
			var root *SingleChildNode[int8]
			var painter *paint.EbitenPainter
			sizeUpdate := true

			props.OnLoad = func(node *scaff.SingleChildNode[*scaff.AcceptNoChild], parent scaff.Node) {

				// Create a single child node that essentially just exists to refresh the builder passed in
				root = &SingleChildNode[int8]{
					id:      "root",
					tracker: node.Tracker(),
					context: &BuildContext{
						cam:          camera,
						BuildContext: node.Context(),
					},
					singleProps: &SingleChildProps[int8]{},
				}
				root.builder = MultiNode(MultiNodeCreate[CanvasProps]{
					ID:           "root-stack",
					PropsCreator: create,
					Create: func(props *MultiChildProps[CanvasProps]) {

						props.OnPropsChanged = func(node *MultiChildNode[CanvasProps]) {
							cv = node.Props()

							// Camera only gets initialized after a while
							if camera != nil {
								// Update camera (size may be changed later)
								camera.LookAt(cv.Position.X, cv.Position.Y)
								camera.SmoothType = cv.SmoothType
								camera.SmoothOptions = &cv.SmoothOptions
							}
						}
					}})
				root.Load(nil)

				cv = root.Children()[0].(*MultiChildNode[CanvasProps]).props
			}

			props.OnUpdate = func(node *scaff.SingleChildNode[*scaff.AcceptNoChild], c *scaff.Context) error {
				if root == nil {
					return nil
				}

				// Forward updates to the child nodes
				return root.Update(c)
			}

			props.OnDraw = func(node *scaff.SingleChildNode[*scaff.AcceptNoChild], c *scaff.Context, image *ebiten.Image) {
				if root == nil {
					return
				}

				// This is run on the first frame as well to create the painter
				if sizeUpdate {
					screen := ebiten.NewImage(image.Bounds().Dx(), image.Bounds().Dy())
					painter = paint.NewEbitenPainter(screen, true, node.Context().AssetManager())

					// Create new camera (old one will only have invalid positions)
					camera = NewCamera(cv.Position.X, cv.Position.Y, c.Width(), c.Height())
					camera.SmoothType = cv.SmoothType
					camera.SmoothOptions = &cv.SmoothOptions
				}

				// Actually draw the root
				painter.Clear()
				painter.SetTransform(paint.Transform{
					CamX:          camera.X,
					CamY:          camera.Y,
					CenterOffsetX: camera.CenterOffsetX,
					CenterOffsetY: camera.CenterOffsetY,
					Angle:         camera.Angle,
					ZoomFactor:    camera.ZoomFactor,
				})
				root.Draw(c, painter)
				image.DrawImage(painter.Screen(), &ebiten.DrawImageOptions{})
			}

			props.OnHandleEvent = func(node *scaff.SingleChildNode[*scaff.AcceptNoChild], c *scaff.Context, event scaff.Event) error {
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
