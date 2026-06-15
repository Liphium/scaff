package engine

import "github.com/Liphium/scaff/scath"

type Texture interface {
	SubTexture(x, y, width, height int) (Texture, error)

	// Origin sets the origin of the next things that are drawn (also for rotation and scaling)
	Origin(origin scath.Vec)

	Rotation(angle float64)

	// Scale scales the next things that are drawn (from Origin)
	Scale(scale scath.Vec)

	Paint(command RenderCommand)

	Clear()
}
