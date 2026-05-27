package uinode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scath"

	"errors"
)

var _ scaff.ChildProps[scaffui.NodeBuilder] = AlignProps{}

type AlignProps struct {
	verticalAlignment  optional.O[VerticalAlignment]
	horizontalAligment optional.O[HorizontalAlignment]

	*scaffui.AcceptChild
}

func (pp *AlignProps) Vertical(alignment VerticalAlignment) {
	pp.verticalAlignment.SetValue(alignment)
}

func (pp *AlignProps) Horizontal(alignment HorizontalAlignment) {
	pp.horizontalAligment.SetValue(alignment)
}

func Align(create func(t *scaff.Tracker, props *AlignProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode(scaffui.SingleNodeCreate[AlignProps]{
		ID: "align",
		DefaultProps: AlignProps{
			horizontalAligment: optional.None[HorizontalAlignment](),
			verticalAlignment:  optional.None[VerticalAlignment](),
			AcceptChild:        &scaffui.AcceptChild{},
		},
		PropsCreator: create,
		Create: func(core *scaffui.SingleChildProps[AlignProps]) {
			// In Layout, we take the biggest we can get in any axis where alignment is given
			core.Layout(func(node *scaffui.SingleChildNode[AlignProps]) (scath.Vec, error) {

				// Pass down constraints from parent to child and let it pick size
				// We just edit this size from now on, since we otherwise want to keep the height / width of our child anyway in case alignment is not set
				size, err := node.LayoutChild(node.Constraints())
				if err != nil {
					return size, err
				}

				if node.Props().horizontalAligment.HasValue() {
					if node.Constraints().MaxX == scath.Infinite {
						return size, errors.New("infinite width for horizontal alignment")
					}

					size.X = node.Constraints().MaxX
				}

				if node.Props().verticalAlignment.HasValue() {
					if node.Constraints().MaxY == scath.Infinite {
						return size, errors.New("infinite height for vertical alignment")
					}

					size.Y = node.Constraints().MaxY
				}

				return size, nil
			})

			// Draw child at proper position for alignment
			core.Draw(func(node *scaffui.SingleChildNode[AlignProps], position scath.Vec, renderer paint.Painter) {
				offset := scath.Vec{}

				if child, ok := node.Child(); ok {
					childSize := child.Current().Size()

					if value, ok := node.Props().horizontalAligment.Value(); ok {
						switch value {
						case HorizontalAlignmentLeft:
							offset.X = 0
						case HorizontalAlignmentCenter:
							offset.X = float64(node.Size().X-childSize.X) / 2
						case HorizontalAlignmentRight:
							offset.X = float64(node.Size().X - childSize.X)
						}
					}

					if value, ok := node.Props().verticalAlignment.Value(); ok {
						switch value {
						case VerticalAlignmentTop:
							offset.Y = 0
						case VerticalAlignmentCenter:
							offset.Y = float64(node.Size().Y-childSize.Y) / 2
						case VerticalAlignmentBottom:
							offset.Y = float64(node.Size().Y - childSize.Y)
						}
					}
				}

				node.DrawChild(position.Add(offset), renderer)
			})
		},
	})
}
