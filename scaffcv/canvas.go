package scaffcv

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scath"
)

type CanvasProps struct {
	size          scath.Vec
	camX, camY    float64
	smoothType    SmoothType
	smoothOptions SmoothOptions
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

func Canvas(create func(t *scaff.Tracker, props *CanvasProps)) scaff.NodeBuilder {
	return scaff.UseNode("canvas", func(props *scaff.SingleChildProps[any]) {
		var cv CanvasProps
		var camera *Camera
		loaded := false

		// Load the props for the node (should be called on state change, etc.)
		fillProps := func(t *scaff.Tracker) {
			cv = CanvasProps{}
			create(t, &cv)
		}

		props.Load(func(node *scaff.SingleChildNode[any], parent scaff.Node) {
			fillProps(node.Tracker())
			loaded = true

			camera = NewCamera(cv.camX, cv.camY, cv.size.X, cv.size.Y)
			node.Tracker()
		})

		props.StateChanged(func(node *scaff.SingleChildNode[any]) {
			// Only do when it isn't the first load
			if !loaded {
				return
			}

			fillProps(node.Tracker())
			camera.SetSize(cv.size.X, cv.size.Y)
			camera.SmoothOptions = &cv.smoothOptions
			camera.SmoothType = cv.smoothType
			camera.LookAt(cv.camX, cv.camY)
		})
	})
}
