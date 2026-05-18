package basenode

import (
	"image/color"

	"github.com/Liphium/scaff/scaffui"
)

type TextProps struct {
	text     string
	fontSize float64
	color    color.RGBA
	*scaffui.AcceptNoChild
}
