package basenode

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
	text           string
	textDirection  text.Direction
	wrapping       bool
	font           string
	fontSize       float64
	lineSpacing    float64
	color          color.Color
	primaryAlign   text.Align
	secondaryAlign text.Align
	*scaffui.AcceptNoChild
}

func (tp *TextProps) Text(text string) {
	tp.text = text
}

func (tp *TextProps) TextDirection(textDirection text.Direction) {
	tp.textDirection = textDirection
}

func (tp *TextProps) Font(font string) {
	tp.font = font
}

func (tp *TextProps) FontSize(size float64) {
	tp.fontSize = size
}

// LineSpacing is a percentage (default 0.25) multiplied with the font size to calculate the spacing between lines. 0 would be no line spacing at all.
func (tp *TextProps) LineSpacing(lineSpacing float64) {
	tp.lineSpacing = lineSpacing
}

func (tp *TextProps) Wrapping(wrapping bool) {
	tp.wrapping = wrapping
}

func (tp *TextProps) Color(color color.Color) {
	tp.color = color
}

// PrimaryAlign sets the primary alignment direction depending on your text direction. If you for example choose left to right as your text direction (the default), this will be horizontal alignment.
func (tp *TextProps) PrimaryAlign(align text.Align) {
	tp.primaryAlign = align
}

// SecondaryAlign sets the secondary alignment direction depending on your text direction. If you for example choose left to right as your text direction (the default), this will be vertical alignment.
func (tp *TextProps) SecondaryAlign(align text.Align) {
	tp.secondaryAlign = align
}

func Text(create func(t *scaff.Tracker, props *TextProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode(scaffui.SingleNodeCreate[TextProps]{
		ID: "text",
		DefaultProps: TextProps{
			text:           "Scaff",
			textDirection:  text.DirectionLeftToRight,
			wrapping:       false,
			fontSize:       16,
			color:          color.White,
			lineSpacing:    0.25,
			primaryAlign:   text.AlignStart,
			secondaryAlign: text.AlignStart,
			AcceptNoChild:  &scaffui.AcceptNoChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.SingleChildProps[TextProps]) {
			finalText := ""

			props.WantedConstraints(func(node *scaffui.SingleChildNode[TextProps], parent scath.Constraints) scath.Constraints {
				font, err := node.Context().AssetManager().GetFont(node.Props().font)
				if err != nil {
					log.Warn("couldn't load font", "font", node.Props().font)
					return scath.Unconstrained()
				}

				// Reusable measure function for these props
				lineSpacing := node.Props().fontSize * node.Props().lineSpacing
				measure := func(t string) (width, height float64) {
					return text.Measure(t, &text.GoTextFace{
						Source:    font,
						Direction: node.Props().textDirection,
						Size:      node.Props().fontSize,
					}, lineSpacing)
				}

				constraints := scath.Unconstrained()
				if !node.Props().wrapping {

					// For no wrapping, just measure the text and set tight constraints, as that's what's wanted
					width, height := measure(node.Props().text)
					constraints = scath.Tight(width, height)
				} else {

					// For no wrapping, we constrain to the longest word + length of one line
					words := sliceIntoWords(node.Props().text)
					w1, h1 := measure(strings.Join(words, "\n"))
					w2, h2 := measure(node.Props().text)
					constraints = scath.NewConstraints(math.Min(w1, w2), math.Max(w1, w2), math.Min(h1, h2), math.Max(h1, h2))
				}

				return constraints
			})

			props.Layout(func(node *scaffui.SingleChildNode[TextProps]) (scath.Vec, error) {
				font, err := node.Context().AssetManager().GetFont(node.Props().font)
				if err != nil {
					log.Warn("couldn't load font", "font", node.Props().font)
					return scath.Vec{}, nil
				}

				measure := func(t string) (width, height float64) {
					return text.Measure(t, &text.GoTextFace{
						Source:    font,
						Direction: node.Props().textDirection,
						Size:      node.Props().fontSize,
					}, 0)
				}

				constraints := node.Constraints()
				maxX := constraints.RealMaxX()
				maxY := constraints.RealMaxY()

				runes := []rune(node.Props().text)
				finalText = ""
				line := ""
				width, height := float64(0), float64(0)
				lineWidth, lineHeight := float64(0), float64(0)
				lineOffset := 0
				lastSpace := 0

				// Line spacing based on the text direction (for width + height)
				lineSpacing := node.Props().fontSize * node.Props().lineSpacing
				lineSpacingWidth, lineSpacingHeight := lineSpacing, lineSpacing
				vertical := node.Props().textDirection == text.DirectionTopToBottomAndLeftToRight || node.Props().textDirection == text.DirectionTopToBottomAndRightToLeft
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
							log.Warn("word is too long for size of text", "w", line)
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

					// If max for a line reached (and the text should wrap), make sure to go back
					if lineWidth > maxX || lineHeight > maxY {
						if node.Props().wrapping {
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
					commitLine(len(node.Props().text)-1, false)
				}

				// Subtract the line spacing again (was added to the last line as well, even though it shouldn't be)
				width -= lineSpacingWidth
				height -= lineSpacingHeight

				return scath.Vec{X: width, Y: height}, nil
			})

			props.Draw(func(node *scaffui.SingleChildNode[TextProps], position scath.Vec, painter paint.Painter) {
				painter.Paint(paint.Text{
					Direction:      node.Props().textDirection,
					Font:           node.Props().font,
					Text:           finalText,
					Color:          node.Props().color,
					FontSize:       node.Props().fontSize,
					LineSpacing:    node.Props().fontSize * node.Props().lineSpacing,
					Position:       position,
					PrimaryAlign:   node.Props().primaryAlign,
					SecondaryAlign: node.Props().secondaryAlign,
				})
			})
		},
	})
}

// This is extracted here so we can improve it in the future, with support for lots of different languages, etc.
func sliceIntoWords(text string) []string {
	return strings.Fields(text)
}
