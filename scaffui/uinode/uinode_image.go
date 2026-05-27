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
	constraints optional.O[scath.Constraints]
	path        optional.O[string]
	filterMode  optional.O[ebiten.Filter]
	*scaffui.AcceptNoChild
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
func (i *ImageProps) Constraints(constraints scath.Constraints) {
	i.constraints.SetValue(constraints)
}

func Image(create func(t *scaff.Tracker, props *ImageProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode(scaffui.SingleNodeCreate[ImageProps]{
		ID: "image",
		DefaultProps: ImageProps{
			AcceptNoChild: &scaffui.AcceptNoChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.SingleChildProps[ImageProps]) {
			props.WantedConstraints(func(node *scaffui.SingleChildNode[ImageProps], parent scath.Constraints) scath.Constraints {
				return node.Props().constraints.Or(scath.Unconstrained())
			})

			props.Layout(func(node *scaffui.SingleChildNode[ImageProps]) (scath.Vec, error) {
				spec := uispec.SingleChildBoxSpec{
					Parent:  node.Constraints(),
					Wanted:  node.Props().constraints,
					Padding: scath.Pad(0),
				}

				return spec.LayoutWithoutChild()
			})

			props.Draw(func(node *scaffui.SingleChildNode[ImageProps], position scath.Vec, painter paint.Painter) {
				if path, ok := node.Props().path.Value(); ok {

					// Draw the actual image
					painter.Paint(paint.Image{
						Path:       path,
						Position:   position,
						Size:       node.Size(),
						FilterMode: node.Props().filterMode.Or(ebiten.FilterLinear),
					})
				} else {

					// Draw a red rectangle to signal an error
					painter.Paint(paint.Rectangle{
						Position:  position,
						Size:      node.Size(),
						FillColor: color.RGBA{255, 0, 0, 255},
					})
				}
			})
		},
	})
}
