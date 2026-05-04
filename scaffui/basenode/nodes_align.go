package basenode

import (
	"errors"

	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/smath"
)

type AlignProps struct {
	child              optional.O[scaffui.NodeBuilder]
	verticalAlignment  optional.O[VerticalAlignment]
	horizontalAligment optional.O[HorizontalAlignment]
}

func (ap *AlignProps) Child(builder scaffui.NodeBuilder) {
	ap.child.SetValue(builder)
}

func (pp *AlignProps) Vertical(alignment VerticalAlignment) {
	pp.verticalAlignment.SetValue(alignment)
}

func (pp *AlignProps) Horizontal(alignment HorizontalAlignment) {
	pp.horizontalAligment.SetValue(alignment)
}

func Align(create func(t *scaffui.Tracker, props *AlignProps)) scaffui.NodeBuilder {
	return scaffui.UseSingleNode("align", create, func(core *scaffui.SingleChildConstruct[AlignProps]) {

		// Pass the child to the core node
		if child, ok := core.Props().child.Value(); ok {
			core.Child(child)
		}

		// In Layout, we take the biggest we can get in any axis where alignment is given
		core.Layout(func(node *scaffui.SingleChildNode[AlignProps]) (scaffui.Size, error) {

			// Pass down constraints from parent to child and let it pick size
			// We just edit this size from now on, since we otherwise want to keep the height / width of our child anyway in case alignment is not set
			size, err := node.LayoutChild(node.Constraints())
			if err != nil {
				return size, err
			}

			if node.Props().horizontalAligment.HasValue() {
				if node.Constraints().MaxWidth == scaffui.Infinite {
					return size, errors.New("infinite width for horizontal alignment")
				}

				size.Width = node.Constraints().MaxWidth
			}

			if node.Props().verticalAlignment.HasValue() {
				if node.Constraints().MaxHeight == scaffui.Infinite {
					return size, errors.New("infinite height for vertical alignment")
				}

				size.Height = node.Constraints().MaxHeight
			}

			return size, nil
		})

		// Draw child at proper position for alignment
		core.Draw(func(node *scaffui.SingleChildNode[AlignProps], position smath.Vec, renderer scaffui.Renderer) {
			offset := smath.Vec{}

			if child, ok := node.Child(); ok {
				childSize := child.Current().Size()

				if horizontal, ok := node.Props().horizontalAligment.Value(); ok {
					switch horizontal {
					case HorizontalAlignmentLeft:
						offset.X = 0
					case HorizontalAlignmentCenter:
						offset.X = float64(node.Size().Width-childSize.Width) / 2
					case HorizontalAlignmentRight:
						offset.X = float64(node.Size().Width - childSize.Width)
					}
				}

				if vertical, ok := node.Props().verticalAlignment.Value(); ok {
					switch vertical {
					case VerticalAlignmentTop:
						offset.Y = 0
					case VerticalAlignmentCenter:
						offset.Y = float64(node.Size().Height-childSize.Height) / 2
					case VerticalAlignmentBottom:
						offset.Y = float64(node.Size().Height - childSize.Height)
					}
				}
			}

			node.DrawChild(position.Add(offset), renderer)
		})
	})
}
