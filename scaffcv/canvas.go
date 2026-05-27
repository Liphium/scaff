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
		props.Load(func(node *scaff.SingleChildNode[any], parent scaff.Node) {
			node.Tracker()
		})
	})
}
