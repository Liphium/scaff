package engine

type Transform struct {
	CamX          float64
	CamY          float64
	CenterOffsetX float64
	CenterOffsetY float64
	Angle         float64 // in radians
	ZoomFactor    float64
}

// NewTransform returns a default Transform with scale (1, 1).
func NewTransform() Transform {
	return Transform{
		ZoomFactor: 1.0,
	}
}
