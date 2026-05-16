package basenode

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
	constraints optional.O[scaffui.Constraints]
	path        optional.O[string]
	filterMode  optional.O[ebiten.Filter]
}

// The path to the image (renderer-specific but probably to the file in your assets file system)
func (i *ImageProps) Path(path string) {
	i.path.SetValue(path)
}

// Set the filter mode used for the image
func (i *ImageProps) Filter(filter ebiten.Filter) {
	i.filterMode.SetValue(filter)
}

// Set the constraints for the image
func (i *ImageProps) Constraints(constraints scaffui.Constraints) {
	i.constraints.SetValue(constraints)
}

func Image(create func(t *scaff.Tracker, props *ImageProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode("image", create, func(core *scaffui.SingleChildProps[ImageProps]) {
		core.WantedConstraints(func(node *scaffui.SingleChildNode[ImageProps], parent scaffui.Constraints) scaffui.Constraints {
			return node.Props().constraints.Or(scaffui.Unconstrained())
		})

		core.Layout(func(node *scaffui.SingleChildNode[ImageProps]) (scath.Vec, error) {
			spec := uispec.SingleChildBoxSpec{
				Parent:        node.Constraints(),
				Wanted:        node.Props().constraints,
				scath.Padding: scaffui.Pad(0),
			}

			return spec.LayoutWithoutChild()
		})

		core.Draw(func(node *scaffui.SingleChildNode[ImageProps], position scath.Vec, renderer paint.Painter) {
			if path, ok := node.Props().path.Value(); ok {

				// Draw the actual image
				renderer.DrawOne(scaffui.ImageCommand{
					Path:       path,
					Position:   position,
					Size:       node.Size(),
					FilterMode: node.Props().filterMode.Or(ebiten.FilterLinear),
				})
			} else {

				// Draw a red rectangle to signal an error
				renderer.DrawOne(scaffui.RectangleCommand{
					Position:  position,
					Size:      node.Size(),
					FillColor: color.RGBA{255, 0, 0, 255},
				})
			}
		})
	})
}
