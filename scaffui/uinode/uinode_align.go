package uinode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/engine"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scath"
)

var _ scaff.ChildProps[scaffui.NodeBuilder] = AlignProps{}

type AlignProps struct {
	VerticalAlignment  optional.O[VerticalAlignment]
	HorizontalAligment optional.O[HorizontalAlignment]

	*scaffui.AcceptChild
}

func Align(create func(t *scaff.Tracker, props *AlignProps)) scaffui.NodeBuilder {
	return scaffui.Standard(scaffui.StandardCreate[AlignProps]{
		ID: "align",
		DefaultProps: AlignProps{
			HorizontalAligment: optional.None[HorizontalAlignment](),
			VerticalAlignment:  optional.None[VerticalAlignment](),
			AcceptChild:        &scaffui.AcceptChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.StandardMethods[AlignProps]) {
			// In Layout, we take the biggest we can get in any axis where alignment is given
			props.OnLayout = func(node *scaffui.StandardNode[AlignProps]) (scath.Vec, error) {

				// Pass down constraints from parent to child and let it pick size
				// We just edit this size from now on, since we otherwise want to keep the height / width of our child anyway in case alignment is not set
				size, err := node.LayoutChildren()
				if err != nil {
					return size, err
				}

				if node.Props().HorizontalAligment.HasValue() && node.Constraints().MaxX != scath.Infinite {
					size.X = node.Constraints().MaxX
				}

				if node.Props().VerticalAlignment.HasValue() && node.Constraints().MaxY != scath.Infinite {
					size.Y = node.Constraints().MaxY
				}

				return size, nil
			}

			// Draw child at proper position for alignment
			props.OnDraw = func(node *scaffui.StandardNode[AlignProps], position scath.Vec, renderer engine.Painter) {
				offset := scath.Vec{}

				if child, ok := node.Child(); ok {
					childSize := child.Size()

					if value, ok := node.Props().HorizontalAligment.Value(); ok {
						switch value {
						case HorizontalAlignmentLeft:
							offset.X = 0
						case HorizontalAlignmentCenter:
							offset.X = float64(node.Size().X-childSize.X) / 2
						case HorizontalAlignmentRight:
							offset.X = float64(node.Size().X - childSize.X)
						}
					}

					if value, ok := node.Props().VerticalAlignment.Value(); ok {
						switch value {
						case VerticalAlignmentTop:
							offset.Y = 0
						case VerticalAlignmentCenter:
							offset.Y = float64(node.Size().Y-childSize.Y) / 2
						case VerticalAlignmentBottom:
							offset.Y = float64(node.Size().Y - childSize.Y)
						}
					}

					child.Draw(position.Add(offset), renderer)
				}
			}
		},
	})
}
