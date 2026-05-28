package scaffcv

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/hajimehoshi/ebiten/v2"
)

type CanvasProps struct {
	camX, camY    float64
	smoothType    SmoothType
	smoothOptions SmoothOptions

	children []NodeBuilder
}

func (cp *CanvasProps) CameraPosition(camX, camY float64) {
	cp.camX = camX
	cp.camY = camY
}

func (cp *CanvasProps) SmoothType(smoothType SmoothType) {
	cp.smoothType = smoothType
}

func (cp *CanvasProps) SmoothOptions(options SmoothOptions) {
	cp.smoothOptions = options
}

func (sp *CanvasProps) Child(builder NodeBuilder) {
	sp.children = append(sp.children, builder)
}

func (sp CanvasProps) GetBuilders() []NodeBuilder {
	return sp.children
}

func Canvas(create func(t *scaff.Tracker, props *CanvasProps)) scaff.NodeBuilder {
	return scaff.UseNode("canvas", func(props *scaff.SingleChildProps[any]) {
		var camera *Camera
		var cv CanvasProps
		var root *SingleChildNode[int8]
		var painter *paint.EbitenPainter
		sizeUpdate := true

		props.Load(func(node *scaff.SingleChildNode[any], parent scaff.Node) {

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
			root.builder = CreateMultiNode("root-stack", create, func(props *MultiChildProps[CanvasProps]) {

				props.PropsChanged(func(node *MultiChildNode[CanvasProps]) {
					cv = node.Props()

					// Camera only gets initialized after a while
					if camera != nil {
						// Update camera (size may be changed later)
						camera.LookAt(cv.camX, cv.camY)
						camera.SmoothType = cv.smoothType
						camera.SmoothOptions = &cv.smoothOptions
					}
				})
			})
			root.Load(nil)

			cv = root.Children()[0].(*MultiChildNode[CanvasProps]).props
		})

		props.Update(func(node *scaff.SingleChildNode[any], c *scaff.Context) error {
			if root == nil {
				return nil
			}

			// Forward updates to the child nodes
			return root.Update(c)
		})

		props.Draw(func(node *scaff.SingleChildNode[any], c *scaff.Context, image *ebiten.Image) {
			if root == nil {
				return
			}

			// This is run on the first frame as well to create the painter
			if sizeUpdate {
				screen := ebiten.NewImage(image.Bounds().Dx(), image.Bounds().Dy())
				painter = paint.NewEbitenPainter(screen, true, node.Context().AssetManager())

				// Create new camera (old one will only have invalid positions)
				camera = NewCamera(cv.camX, cv.camY, c.Width(), c.Height())
				camera.SmoothType = cv.smoothType
				camera.SmoothOptions = &cv.smoothOptions
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
		})

		props.HandleEvent(func(node *scaff.SingleChildNode[any], c *scaff.Context, event scaff.Event) error {
			if root == nil {
				return nil
			}

			// For a size change, make sure we properly handle it
			if event.EventID() == scaff.EventIdSizeChange {
				sizeUpdate = true
			}

			// Forward events to the root
			return root.HandleEvent(c, event)
		})
	})
}
