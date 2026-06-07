package uinode

import (
	"image"

	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffcv/cvnode"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uispec"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type ImageProps struct {
	Constraints optional.O[scath.Constraints]
	Path        string
	FilterMode  ebiten.Filter
	Offset      optional.O[cvnode.TilePosition]
	Frame       optional.O[cvnode.TilePosition]
	*scaffui.AcceptNoChild
}

// The Path to the image (renderer-specific but probably to the file in your assets file system)
// Set the filter mode used for the image
// Set the Constraints for the image
func Image(create func(t *scaff.Tracker, props *ImageProps)) scaffui.NodeBuilder {
	return scaffui.Standard(scaffui.StandardCreate[ImageProps]{
		ID: "image",
		DefaultProps: ImageProps{
			FilterMode:    ebiten.FilterLinear,
			Offset:        optional.None[cvnode.TilePosition](),
			Frame:         optional.None[cvnode.TilePosition](),
			AcceptNoChild: &scaffui.AcceptNoChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.StandardMethods[ImageProps]) {
			props.OnWantedConstraints = func(node *scaffui.StandardNode[ImageProps], parent scath.Constraints) scath.Constraints {
				return node.Props().Constraints.Or(scath.Unconstrained())
			}

			props.OnLayout = func(node *scaffui.StandardNode[ImageProps]) (scath.Vec, error) {
				spec := uispec.SingleChildBoxSpec{
					Parent:  node.Constraints(),
					Wanted:  node.Props().Constraints,
					Padding: scath.Pad(0),
				}

				return spec.LayoutWithoutChild()
			}

			props.OnDraw = func(node *scaffui.StandardNode[ImageProps], position scath.Vec, painter paint.Painter) {
				asset, err := node.Context().AssetManager().GetImage(node.Props().Path)
				if err != nil {
					log.Error("couldn't find image", "i", node.Props().Path)
					return
				}

				// If it is a sub-image make sure to cut it out
				width, height := 0.0, 0.0
				if offset, ok := node.Props().Offset.Value(); ok {
					frame, ok := node.Props().Frame.Value()
					if !ok {
						log.Error("can't render subimage without frame")
						return
					}

					asset = asset.SubImage(image.Rect(offset.X, offset.Y, offset.X+frame.X, offset.Y+frame.Y)).(*ebiten.Image)
					width, height = float64(frame.X), float64(frame.Y)
				} else {
					width, height = float64(asset.Bounds().Dx()), float64(asset.Bounds().Dy())
				}

				op := &ebiten.DrawImageOptions{
					Filter: node.Props().FilterMode,
				}
				realSize := node.Size()
				if realSize.X != 0 && realSize.Y != 0 {
					op.GeoM.Scale(realSize.X/width, realSize.Y/height)
				}
				op.GeoM.Translate(position.X, position.Y)

				// Draw the actual image
				painter.DrawRaw(asset, op)
			}
		},
	})
}
