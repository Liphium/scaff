package basenode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/smath"
)

type FlexProps struct {
	children      []scaffui.NodeBuilder
	stretchFactor map[int]optional.O[int]
	direction     optional.O[LayoutDirection]
}

func (fp *FlexProps) Child(builder scaffui.NodeBuilder) {
	if fp.stretchFactor == nil {
		fp.stretchFactor = map[int]optional.O[int]{}
	}

	fp.stretchFactor[len(fp.children)] = optional.None[int]()
	fp.children = append(fp.children, builder)
}

func (fp *FlexProps) Expanded(factor int, builder scaffui.NodeBuilder) {
	if fp.stretchFactor == nil {
		fp.stretchFactor = map[int]optional.O[int]{}
	}

	fp.stretchFactor[len(fp.children)] = optional.With(factor)
	fp.children = append(fp.children, builder)
}

func (fp *FlexProps) Direction(direction LayoutDirection) {
	fp.direction.SetValue(direction)
}

func Flex(create func(t *scaff.Tracker, props *FlexProps)) scaffui.NodeBuilder {
	return scaffui.UseMultiNode("flex", create, func(core *scaffui.MultiChildConstruct[FlexProps]) {

		// Pass all of the children to the multi child node
		for _, child := range core.Props().children {
			core.Child(child)
		}

		core.Layout(func(node *scaffui.MultiChildNode[FlexProps]) (scaffui.Size, error) {
			return flexLayout(node)
		})

		core.Draw(func(node *scaffui.MultiChildNode[FlexProps], position smath.Vec, renderer scaffui.Renderer) {
			flexDraw(node, position, renderer)
		})
	})
}

func flexDirection(node *scaffui.MultiChildNode[FlexProps]) LayoutDirection {
	return node.Props().direction.Or(LayoutTopToBottom)
}

func flexDraw(node *scaffui.MultiChildNode[FlexProps], position smath.Vec, renderer scaffui.Renderer) {
	children := node.Children()

	switch flexDirection(node) {
	case LayoutRightToLeft:
		xOffset := float64(node.Size().Width)
		for _, child := range children {
			childSize := child.Current().Size()
			xOffset -= float64(childSize.Width)
			child.Current().Draw(smath.Vec{X: position.X + xOffset, Y: position.Y}, renderer)
		}

	case LayoutTopToBottom:
		yOffset := 0.0
		for _, child := range children {
			child.Current().Draw(smath.Vec{X: position.X, Y: position.Y + yOffset}, renderer)
			yOffset += float64(child.Current().Size().Height)
		}

	case LayoutBottomToTop:
		yOffset := float64(node.Size().Height)
		for _, child := range children {
			childSize := child.Current().Size()
			yOffset -= float64(childSize.Height)
			child.Current().Draw(smath.Vec{X: position.X, Y: position.Y + yOffset}, renderer)
		}

	case LayoutLeftToRight:
		xOffset := 0.0
		for _, child := range children {
			child.Current().Draw(smath.Vec{X: position.X + xOffset, Y: position.Y}, renderer)
			xOffset += float64(child.Current().Size().Width)
		}

	default:
		log.Warn("invalid layout direction on flex layout node")
	}
}

func flexLayout(node *scaffui.MultiChildNode[FlexProps]) (scaffui.Size, error) {
	switch flexDirection(node) {
	case LayoutBottomToTop:
		fallthrough // Same as top to bottom cause no positions considered yet
	case LayoutTopToBottom:
		return flexLayoutLinear(node, false)
	case LayoutRightToLeft:
		fallthrough // Same as left to right cause no positions considered yet
	case LayoutLeftToRight:
		return flexLayoutLinear(node, true)
	}

	return scaffui.Size{}, nil
}

// Shared linear layout logic for horizontal and vertical directions.
func flexLayoutLinear(node *scaffui.MultiChildNode[FlexProps], horizontal bool) (scaffui.Size, error) {
	mainMin, mainMax, crossMin, crossMax := axisConstraints(node.Constraints(), horizontal)
	children := node.Children()
	totalSize := scaffui.Size{}

	if len(children) == 0 {
		return clampLinearSize(totalSize, horizontal, mainMin, mainMax, crossMin, crossMax), nil
	}

	if mainMax == scaffui.Infinite {
		for _, child := range children {
			child.Current().SetConstraints(newChildConstraints(horizontal, 0, scaffui.Infinite, crossMin, crossMax))
			size, err := child.Current().Layout()
			if err != nil {
				return scaffui.Size{}, err
			}
			totalSize = addChildSize(totalSize, size, horizontal)
		}

		return clampLinearSize(totalSize, horizontal, mainMin, mainMax, crossMin, crossMax), nil
	}

	remainder := mainMax
	childCount := len(children)
	share := remainder / childCount

	// Lay out non-stretchable children
	for i, child := range node.Children() {
		factor, ok := node.Props().stretchFactor[i]
		if ok && factor.HasValue() {
			continue
		}

		// Lay out the child
		child.Current().SetConstraints(newChildConstraints(horizontal, 0, share, crossMin, crossMax))
		size, err := child.Current().Layout()
		if err != nil {
			return scaffui.Size{}, err
		}
		totalSize = addChildSize(totalSize, size, horizontal)

		remainder -= mainSize(size, horizontal)
		childCount--
		if childCount > 0 {
			share = remainder / childCount
		}
	}

	// Add all the factors together
	factorSum := 0
	for i := range children {
		factor, ok := node.Props().stretchFactor[i]
		if !ok || !factor.HasValue() {
			continue
		}

		factorSum += factor.Or(0)
	}

	// Lay out stretchable children
	for i, child := range children {
		factor, ok := node.Props().stretchFactor[i]
		if !ok || !factor.HasValue() {
			continue
		}

		child.Current().SetConstraints(newChildConstraints(horizontal, 0, remainder/factorSum*factor.Or(1), crossMin, crossMax))
		size, err := child.Current().Layout()
		if err != nil {
			return scaffui.Size{}, err
		}
		totalSize = addChildSize(totalSize, size, horizontal)
	}

	return clampLinearSize(totalSize, horizontal, mainMin, mainMax, crossMin, crossMax), nil
}

// Maps constraints to main/cross axis values based on direction.
func axisConstraints(c scaffui.Constraints, horizontal bool) (mainMin, mainMax, crossMin, crossMax int) {
	if horizontal {
		return c.MinWidth, c.MaxWidth, c.MinHeight, c.MaxHeight
	}

	return c.MinHeight, c.MaxHeight, c.MinWidth, c.MaxWidth
}

// Builds child constraints from main/cross axis values.
func newChildConstraints(horizontal bool, mainMin, mainMax, crossMin, crossMax int) scaffui.Constraints {
	if horizontal {
		return scaffui.NewConstraints(mainMin, mainMax, crossMin, crossMax)
	}

	return scaffui.NewConstraints(crossMin, crossMax, mainMin, mainMax)
}

// Find the max size in the axis of a child.
func mainSize(size scaffui.Size, horizontal bool) int {
	if horizontal {
		return size.Width
	}

	return size.Height
}

// Aggregates child size into total container size.
func addChildSize(total, child scaffui.Size, horizontal bool) scaffui.Size {
	if horizontal {
		total.Width += child.Width
		if child.Height > total.Height {
			total.Height = child.Height
		}
		return total
	}

	total.Height += child.Height
	if child.Width > total.Width {
		total.Width = child.Width
	}
	return total
}

func clampLinearSize(size scaffui.Size, horizontal bool, mainMin, mainMax, crossMin, crossMax int) scaffui.Size {
	main := mainSize(size, horizontal)
	cross := mainSize(size, !horizontal)

	main = clampAxis(main, mainMin, mainMax)
	cross = clampAxis(cross, crossMin, crossMax)

	if horizontal {
		return scaffui.Size{Width: main, Height: cross}
	}

	return scaffui.Size{Width: cross, Height: main}
}

func clampAxis(v, minV, maxV int) int {
	v = max(v, minV)
	if maxV != scaffui.Infinite {
		v = min(v, maxV)
	}
	return v
}
