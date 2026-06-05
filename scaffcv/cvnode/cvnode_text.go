package cvnode

import (
	"image/color"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TextProps struct {
	// Path to the font inside of the assets folder
	Font string

	// Font size of the text
	FontSize float64

	// LineSpacing is a percentage (default 0.25) multiplied with the Font size to calculate the spacing between lines. 0 would be no line spacing at all.
	LineSpacing float64

	// Actual text that should be displayed (can contain \n for new lines)
	Text string

	// Direction of the text
	TextDirection text.Direction

	// Position of the text
	Position scath.Vec

	// Color of the text
	Color color.Color

	// Primary alignment of the text (for example: horizontal alignment when left to right / right to left)
	PrimaryAlign text.Align

	// Secondary alignment of the text (for example: vertical alignment when left to right / right to left)
	SecondaryAlign text.Align

	scaffcv.AcceptNoChild
}

// Returns topLeft position for the size of the text
func (t TextProps) textPosition(size scath.Vec) scath.Vec {
	position := t.Position

	// Calculate text direction into the position
	multiplier := -1.0
	switch t.TextDirection {
	case text.DirectionTopToBottomAndRightToLeft, text.DirectionRightToLeft:
		position = position.Add(scath.Vec{X: -size.X})
		multiplier = 1
	}

	horizontal := t.TextDirection == text.DirectionLeftToRight || t.TextDirection == text.DirectionTopToBottomAndLeftToRight

	getForAlignment := func(alignment text.Align, value float64) float64 {
		switch alignment {
		case text.AlignCenter:
			return multiplier * value / 2
		case text.AlignEnd:
			return multiplier * value
		}
		return 0
	}

	// Calculate primary alignment into the position
	if horizontal {
		position = position.Add(scath.Vec{X: getForAlignment(t.PrimaryAlign, size.X)})
	} else {
		position = position.Add(scath.Vec{Y: getForAlignment(t.PrimaryAlign, size.Y)})
	}

	// Calculate secondary alignment into the position
	if horizontal {
		position = position.Add(scath.Vec{Y: getForAlignment(t.SecondaryAlign, size.Y)})
	} else {
		position = position.Add(scath.Vec{X: getForAlignment(t.SecondaryAlign, size.X)})
	}

	return position
}

// Text creates a simple Text node with a Position and Text.
func Text(create func(t *scaff.Tracker, props *TextProps)) scaffcv.NodeBuilder {
	return scaffcv.Standard(scaffcv.StandardCreate[TextProps]{
		ID: "Text",
		DefaultProps: TextProps{
			Text:           "",
			TextDirection:  text.DirectionLeftToRight,
			LineSpacing:    0.25,
			FontSize:       24,
			Position:       scath.Vec{},
			Color:          color.White,
			PrimaryAlign:   text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
		},
		PropsCreator: create,
		Create: func(props *scaffcv.StandardMethods[TextProps]) {
			size := scath.Zero

			// Recalculate size on every props change
			props.OnPropsChanged = func(node *scaffcv.StandardNode[TextProps]) {
				font, err := node.Context().AssetManager().GetFont(node.Props().Font)
				if err != nil {
					log.Error("couldn't load font for text measuring", "f", node.Props().Font)
					return
				}

				lineSpacing := node.Props().FontSize * node.Props().LineSpacing
				size.X, size.Y = text.Measure(node.Props().Text, &text.GoTextFace{
					Source:    font,
					Size:      node.Props().FontSize,
					Direction: node.Props().TextDirection,
				}, lineSpacing)
			}

			props.Position = func(node *scaffcv.StandardNode[TextProps]) scath.Vec {
				return node.Props().textPosition(size)
			}
			props.Size = func(node *scaffcv.StandardNode[TextProps]) scath.Vec {
				return size
			}

			props.OnDraw = func(node *scaffcv.StandardNode[TextProps], c *scaff.Context, painter paint.Painter) {
				lineSpacing := node.Props().FontSize * node.Props().LineSpacing
				painter.Paint(paint.Text{
					Font:           node.Props().Font,
					Text:           node.Props().Text,
					Direction:      node.Props().TextDirection,
					LineSpacing:    lineSpacing,
					FontSize:       node.Props().FontSize,
					Position:       node.Props().Position,
					Color:          node.Props().Color,
					PrimaryAlign:   node.Props().PrimaryAlign,
					SecondaryAlign: node.Props().SecondaryAlign,
				})
			}
		},
	})
}
