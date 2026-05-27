package scaffcv

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type CanvasProps struct {
	size          scath.Vec
	camX, camY    float64
	smoothType    SmoothType
	smoothOptions SmoothOptions

	children []NodeBuilder
}

func (cp *CanvasProps) Size(size scath.Vec) {
	cp.size = size
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
		loaded := false

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
			root.builder = CreateMultiNode("root-stack", create, nil)
			root.Load(nil)
			loaded = true

			cv = root.Children()[0].(*MultiChildNode[CanvasProps]).props
			camera = NewCamera(cv.camX, cv.camY, cv.size.X, cv.size.Y)
			node.Tracker()
		})

		props.PropsChanged(func(node *scaff.SingleChildNode[any]) {
			// Only do when it isn't the first load
			if !loaded {
				return
			}

			// Synchronize the tracker + props again
			root.tracker = node.Tracker()
			cv = root.Children()[0].(*MultiChildNode[CanvasProps]).props

			// Update the camera properly based on the props
			camera.SetSize(cv.size.X, cv.size.Y)
			camera.SmoothOptions = &cv.smoothOptions
			camera.SmoothType = cv.smoothType
			camera.LookAt(cv.camX, cv.camY)
		})

		props.Draw(func(node *scaff.SingleChildNode[any], c *scaff.Context, image *ebiten.Image) {

		})
	})
}
