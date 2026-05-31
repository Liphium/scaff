package uinode

import (
	"github.com/Liphium/scaff/paint"

	"image/color"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uispec"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type ImageProps struct {
	Constraints optional.O[scath.Constraints]
	Path        optional.O[string]
	FilterMode  optional.O[ebiten.Filter]
	*scaffui.AcceptNoChild
}

// The Path to the image (renderer-specific but probably to the file in your assets file system)
// Set the filter mode used for the image
// Set the Constraints for the image
func Image(create func(t *scaff.Tracker, props *ImageProps)) scaffui.NodeBuilder {
	return scaffui.SingleNode(scaffui.SingleNodeCreate[ImageProps]{
		ID: "image",
		DefaultProps: ImageProps{
			AcceptNoChild: &scaffui.AcceptNoChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.SingleChildProps[ImageProps]) {
			props.OnWantedConstraints = func(node *scaffui.SingleChildNode[ImageProps], parent scath.Constraints) scath.Constraints {
				return node.Props().Constraints.Or(scath.Unconstrained())
			}

			props.OnLayout = func(node *scaffui.SingleChildNode[ImageProps]) (scath.Vec, error) {
				spec := uispec.SingleChildBoxSpec{
					Parent:  node.Constraints(),
					Wanted:  node.Props().Constraints,
					Padding: scath.Pad(0),
				}

				return spec.LayoutWithoutChild()
			}

			props.OnDraw = func(node *scaffui.SingleChildNode[ImageProps], position scath.Vec, painter paint.Painter) {
				if Path, ok := node.Props().Path.Value(); ok {

					// Draw the actual image
					painter.Paint(paint.Image{
						Path:       Path,
						Position:   position,
						Size:       node.Size(),
						FilterMode: node.Props().FilterMode.Or(ebiten.FilterLinear),
					})
				} else {

					// Draw a red rectangle to signal an error
					painter.Paint(paint.Rectangle{
						Position:  position,
						Size:      node.Size(),
						FillColor: color.RGBA{255, 0, 0, 255},
					})
				}
			}
		},
	})
}
