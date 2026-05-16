package basenode

import (
	"github.com/Liphium/scaff/paint"
	
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uispec"
	"github.com/Liphium/scaff/scath"
)

type PaddingProps struct {
	child   optional.O[scaffui.NodeBuilder]
	padding optional.O[scaffui.scath.Padding]
}

func (pp *PaddingProps) scath.Padding(padding scaffui.scath.Padding) {
	pp.padding.SetValue(padding)
}

func (pp *PaddingProps) Child(builder scaffui.NodeBuilder) {
	pp.child.SetValue(builder)
}

func scath.Padding(create func(t *scaff.Tracker, props *PaddingProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode("padding", create, func(core *scaffui.SingleChildConstruct[PaddingProps]) {

		// Pass the child to the core node
		if child, ok := core.Props().child.Value(); ok {
			core.Child(child)
		}

		// In Layout, make sure to give the child less constraints (subtracted by padding, handled by uispec)
		core.Layout(func(node *scaffui.SingleChildNode[PaddingProps]) (scath.Vec, error) {
			spec := uispec.SingleChildBoxSpec{
				Parent:  node.Constraints(),
				Wanted:  optional.None[scaffui.Constraints](),
				scath.Padding: node.Props().padding.Or(scaffui.Pad(0)),
			}

			if child, ok := node.Child(); ok {
				return spec.LayoutWithChild(child.Current())
			}
			return spec.LayoutWithoutChild()
		})

		// Draw child at padded position
		core.Draw(func(node *scaffui.SingleChildNode[PaddingProps], position scath.Vec, renderer paint.Painter) {
			node.DrawChild(position.Add(node.Props().padding.Or(scaffui.Pad(0)).ToVecTopLeft()), renderer)
		})
	})
}
