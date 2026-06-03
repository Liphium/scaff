package uinode

import (
	"image/color"
	"math"
	"strings"
	"unicode"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TextProps struct {
	Text          string
	TextDirection text.Direction
	Wrapping      bool
	Font          string
	FontSize      float64

	// LineSpacing is a percentage (default 0.25) multiplied with the Font size to calculate the spacing between lines. 0 would be no line spacing at all.
	LineSpacing float64

	Color color.Color

	// PrimaryAlign sets the primary alignment direction depending on your Text direction. If you for example choose left to right as your Text direction (the default), this will be horizontal alignment.
	PrimaryAlign text.Align

	// SecondaryAlign sets the secondary alignment direction depending on your Text direction. If you for example choose left to right as your Text direction (the default), this will be vertical alignment.
	SecondaryAlign text.Align

	*scaffui.AcceptNoChild
}

func Text(create func(t *scaff.Tracker, props *TextProps)) scaffui.NodeBuilder {
	return scaffui.Standard(scaffui.StandardCreate[TextProps]{
		ID: "Text",
		DefaultProps: TextProps{
			Text:           "Scaff",
			TextDirection:  text.DirectionLeftToRight,
			Wrapping:       false,
			FontSize:       16,
			Color:          color.White,
			LineSpacing:    0.25,
			PrimaryAlign:   text.AlignStart,
			SecondaryAlign: text.AlignStart,
			AcceptNoChild:  &scaffui.AcceptNoChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.StandardMethods[TextProps]) {
			finalText := ""

			props.OnWantedConstraints = func(node *scaffui.StandardNode[TextProps], parent scath.Constraints) scath.Constraints {
				font, err := node.Context().AssetManager().GetFont(node.Props().Font)
				if err != nil {
					log.Warn("couldn't load font", "font", node.Props().Font)
					return scath.Unconstrained()
				}

				// Reusable measure function for these props
				lineSpacing := node.Props().FontSize * node.Props().LineSpacing
				measure := func(t string) (width, height float64) {
					return text.Measure(t, &text.GoTextFace{
						Source:    font,
						Direction: node.Props().TextDirection,
						Size:      node.Props().FontSize,
					}, lineSpacing)
				}

				constraints := scath.Unconstrained()
				if !node.Props().Wrapping {

					// For no wrapping, just measure the text and set tight constraints, as that's what's wanted
					width, height := measure(node.Props().Text)
					constraints = scath.Tight(width, height)
				} else {

					// For no wrapping, we constrain to the longest word + length of one line
					words := sliceIntoWords(node.Props().Text)
					w1, h1 := measure(strings.Join(words, "\n"))
					w2, h2 := measure(node.Props().Text)
					constraints = scath.NewConstraints(math.Min(w1, w2), math.Max(w1, w2), math.Min(h1, h2), math.Max(h1, h2))
				}

				return constraints
			}

			props.OnLayout = func(node *scaffui.StandardNode[TextProps]) (scath.Vec, error) {
				font, err := node.Context().AssetManager().GetFont(node.Props().Font)
				if err != nil {
					log.Warn("couldn't load Font", "Font", node.Props().Font)
					return scath.Vec{}, nil
				}

				measure := func(t string) (width, height float64) {
					return text.Measure(t, &text.GoTextFace{
						Source:    font,
						Direction: node.Props().TextDirection,
						Size:      node.Props().FontSize,
					}, 0)
				}

				constraints := node.Constraints()
				maxX := constraints.RealMaxX()
				maxY := constraints.RealMaxY()

				runes := []rune(node.Props().Text)
				finalText = ""
				line := ""
				width, height := float64(0), float64(0)
				lineWidth, lineHeight := float64(0), float64(0)
				lineOffset := 0
				lastSpace := 0

				// Line spacing based on the text direction (for width + height)
				lineSpacing := node.Props().FontSize * node.Props().LineSpacing
				lineSpacingWidth, lineSpacingHeight := lineSpacing, lineSpacing
				vertical := node.Props().TextDirection == text.DirectionTopToBottomAndLeftToRight || node.Props().TextDirection == text.DirectionTopToBottomAndRightToLeft
				if vertical {
					lineSpacingHeight = 0
				} else {
					lineSpacingWidth = 0
				}

				// Updates a lines measurement properly
				updateLineMeasurement := func() {
					lineWidth, lineHeight = measure(line)
					lineWidth, lineHeight = lineWidth+lineSpacingWidth, lineHeight+lineSpacingHeight
				}

				// Commits a line and returns new index
				commitLine := func(i int, tillSpace bool) int {
					if line == "" {
						return i
					}

					// Cut line to stuff with last space
					if tillSpace {
						if lastSpace <= lineOffset {
							// There is no space within the current line, cut the line till the next space (this is just to clip text, but is in fact a rendering error)
							log.Warn("word is too long for size of Text", "w", line)
							line = line[0:i]
							found := false
							for j, rune := range runes[i:] {
								if unicode.IsSpace(rune) {
									lineOffset = i + j + 1
									found = true
									break
								}
							}
							if !found {
								lineOffset = len(runes)
							}
						} else {
							// There was a space, wrap after that space
							line = line[0:(lastSpace - lineOffset)]
							lineOffset = lastSpace + 1
						}
					} else {
						lineOffset = i
					}

					// Measure line again: needs to be done due to the last char being calculated in when cutting off text
					updateLineMeasurement()

					// Add to global width / height
					if vertical {
						width += lineWidth
						height = math.Max(height, lineHeight)
					} else {
						width = math.Max(width, lineWidth)
						height += lineHeight
					}

					finalText += line + "\n"
					line = ""

					return lineOffset
				}

				i := 0
				for i < len(runes) {
					rune := runes[i]
					char := string(rune)
					line += char
					updateLineMeasurement()

					// If the global limit is ever reached, just stop calculating
					if width+lineWidth >= maxX && height+lineHeight >= maxY {
						i = commitLine(i, false)
						break
					}

					// If max for a line reached (and the Text should wrap), make sure to go back
					if lineWidth > maxX || lineHeight > maxY {
						if node.Props().Wrapping {
							i = commitLine(i, true)
							continue
						}

						i = commitLine(i, false)
						break
					}

					if unicode.IsSpace(rune) {
						lastSpace = i
					}

					i++
				}

				// Commit the line when not committed yet
				if line != "" {
					commitLine(len(node.Props().Text)-1, false)
				}

				// Subtract the line spacing again (was added to the last line as well, even though it shouldn't be)
				width -= lineSpacingWidth
				height -= lineSpacingHeight

				return scath.Vec{X: width, Y: height}, nil
			}

			props.OnDraw = func(node *scaffui.StandardNode[TextProps], position scath.Vec, painter paint.Painter) {
				painter.Paint(paint.Text{
					Direction:      node.Props().TextDirection,
					Font:           node.Props().Font,
					Text:           finalText,
					Color:          node.Props().Color,
					FontSize:       node.Props().FontSize,
					LineSpacing:    node.Props().FontSize * node.Props().LineSpacing,
					Position:       position,
					PrimaryAlign:   node.Props().PrimaryAlign,
					SecondaryAlign: node.Props().SecondaryAlign,
				})
			}
		},
	})
}

// This is extracted here so we can improve it in the future, with support for lots of different languages, etc.
func sliceIntoWords(Text string) []string {
	return strings.Fields(Text)
}
