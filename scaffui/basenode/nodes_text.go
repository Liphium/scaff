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
	text          string
	textDirection text.Direction
	wrapping      bool
	font          string
	fontSize      float64
	lineSpacing   float64
	color         color.Color
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

func Text(create func(t *scaff.Tracker, props *TextProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode(scaffui.SingleNodeCreate[TextProps]{
		ID: "text",
		DefaultProps: TextProps{
			text:          "Scaff",
			textDirection: text.DirectionLeftToRight,
			wrapping:      false,
			fontSize:      16,
			color:         color.White,
			lineSpacing:   0.25,
			AcceptNoChild: &scaffui.AcceptNoChild{},
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
							for si, rune := range runes[i:] {
								log.Debug("trying to find space", "r", string(rune))
								if unicode.IsSpace(rune) {
									lineOffset = i + si + 1
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
					log.Debug("line", "l", line)
					finalText += line + "\n"
					line = ""

					// Add to global width / height
					width += lineWidth
					height += lineHeight
					return lineOffset
				}

				i := 0
				for i < len(runes) {
					rune := runes[i]
					char := string(rune)
					log.Debug("iteration", "c", char)
					line += char
					lineWidth, lineHeight = measure(line)
					lineWidth, lineHeight = lineWidth+lineSpacingWidth, lineHeight+lineSpacingHeight

					// If the global limit is ever reached, just stop calculating
					if width+lineWidth >= maxX && height+lineHeight >= maxY {
						i = commitLine(i, false)
						log.Debug("commit, max reached")
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

				return scath.Vec{X: width, Y: height}, nil
			})

			props.Draw(func(node *scaffui.SingleChildNode[TextProps], position scath.Vec, painter paint.Painter) {
				painter.Paint(paint.Text{
					Direction:   node.Props().textDirection,
					Font:        node.Props().font,
					Text:        finalText,
					Color:       node.Props().color,
					FontSize:    node.Props().fontSize,
					LineSpacing: node.Props().fontSize * node.Props().lineSpacing,
					Position:    position,
				})
			})
		},
	})
}

// This is extracted here so we can improve it in the future, with support for lots of different languages, etc.
func sliceIntoWords(text string) []string {
	return strings.Fields(text)
}
