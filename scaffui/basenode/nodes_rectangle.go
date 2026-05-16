package basenode

import (
	"github.com/Liphium/scaff/paint"

	"image/color"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uispec"
	"github.com/Liphium/scaff/scath"
)

type RectangleProps struct {
	child             optional.O[scaffui.NodeBuilder]
	wantedConstraints optional.O[scaffui.Constraints]
	padding           optional.O[scaffui.scath.Padding]
	fillColor         optional.O[color.RGBA]
	borderRadius      optional.O[int]
}

func (rp *RectangleProps) Child(builder scaffui.NodeBuilder) {
	rp.child.SetValue(builder)
}

func (rp *RectangleProps) WantedConstraints(constraints scaffui.Constraints) {
	rp.wantedConstraints.SetValue(constraints)
}

func (rp *RectangleProps) scath.Padding(padding scaffui.scath.Padding) {
	rp.padding.SetValue(padding)
}

func (rp *RectangleProps) FillColor(fillColor color.RGBA) {
	rp.fillColor.SetValue(fillColor)
}

func (rp *RectangleProps) BorderRadius(borderRadius int) {
	rp.borderRadius.SetValue(borderRadius)
}

func Rectangle(create func(t *scaff.Tracker, props *RectangleProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode("rectangle", create, func(core *scaffui.SingleChildProps[RectangleProps]) {

		// Pass the child to the single node (when there)
		if child, ok := core.Props().child.Value(); ok {
			core.Child(child)
		}

		core.WantedConstraints(func(node *scaffui.SingleChildNode[RectangleProps], _ scaffui.Constraints) scaffui.Constraints {
			return core.Props().wantedConstraints.Or(scaffui.Unconstrained())
		})

		core.Layout(func(node *scaffui.SingleChildNode[RectangleProps]) (scath.Vec, error) {
			props := core.Props()
			spec := uispec.SingleChildBoxSpec{
				Parent:  node.Constraints(),
				Wanted:  props.wantedConstraints,
				scath.Padding: props.padding.Or(scaffui.Pad(0)),
			}

			child, ok := node.Child()
			if ok {
				return spec.LayoutWithChild(child.Current())
			}

			return spec.LayoutWithoutChild()
		})

		core.Draw(func(node *scaffui.SingleChildNode[RectangleProps], position scath.Vec, renderer paint.Painter) {
			props := node.Props()
			renderer.DrawOne(scaffui.RectangleCommand{
				Position:     position,
				Size:         node.Size(),
				FillColor:    props.fillColor.Or(color.RGBA{255, 255, 255, 255}),
				BorderRadius: props.borderRadius.Or(0),
			})

			node.DrawChild(position.Add(props.padding.Or(scaffui.Pad(0)).ToVecTopLeft()), renderer)
		})
	})
}
