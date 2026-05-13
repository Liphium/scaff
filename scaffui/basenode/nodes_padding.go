package basenode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uispec"
	"github.com/Liphium/scaff/smath"
)

type PaddingProps struct {
	child   optional.O[scaffui.NodeBuilder]
	padding optional.O[scaffui.Padding]
}

func (pp *PaddingProps) Padding(padding scaffui.Padding) {
	pp.padding.SetValue(padding)
}

func (pp *PaddingProps) Child(builder scaffui.NodeBuilder) {
	pp.child.SetValue(builder)
}

func Padding(create func(t *scaff.Tracker, props *PaddingProps)) scaffui.NodeBuilder {
	return scaffui.UseSingleNode("padding", create, func(core *scaffui.SingleChildConstruct[PaddingProps]) {

		// Pass the child to the core node
		if child, ok := core.Props().child.Value(); ok {
			core.Child(child)
		}

		// In Layout, make sure to give the child less constraints (subtracted by padding, handled by uispec)
		core.Layout(func(node *scaffui.SingleChildNode[PaddingProps]) (scaffui.Size, error) {
			spec := uispec.SingleChildBoxSpec{
				Parent:  node.Constraints(),
				Wanted:  optional.None[scaffui.Constraints](),
				Padding: node.Props().padding.Or(scaffui.Pad(0)),
			}

			if child, ok := node.Child(); ok {
				return spec.LayoutWithChild(child.Current())
			}
			return spec.LayoutWithoutChild()
		})

		// Draw child at padded position
		core.Draw(func(node *scaffui.SingleChildNode[PaddingProps], position smath.Vec, renderer scaffui.Renderer) {
			node.DrawChild(position.Add(node.Props().padding.Or(scaffui.Pad(0)).ToVecTopLeft()), renderer)
		})
	})
}
