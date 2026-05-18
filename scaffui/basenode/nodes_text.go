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

func (tp *TextProps) Color(color color.RGBA) {
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

				lineSpacing := node.Props().fontSize * node.Props().lineSpacing
				measure := func(t string) (width, height float64) {
					return text.Measure(t, &text.GoTextFace{
						Source:    font,
						Direction: node.Props().textDirection,
						Size:      node.Props().fontSize,
					}, lineSpacing)
				}

				constraints := node.Constraints()
				maxX := constraints.RealMaxX()
				maxY := constraints.RealMaxY()

				text := ""
				line := ""
				width, height := float64(0), float64(0)
				lineWidth, lineHeight := float64(0), float64(0)
				lineOffset := 0
				lastSpace := 0

				// TODO:
				// - Properly handle line spacing (based on text direction)
				// - Properly calculate global size

				i := 0
				for i < len(node.Props().text) {
					rune := rune(node.Props().text[i])
					char := string(rune)
					w, h := measure(char)
					lineWidth, lineHeight = lineWidth+w, lineHeight+h

					// If the global limit is ever reached, just stop calculating
					if width+lineWidth > maxX || height+lineHeight > maxY {
						// Cut line to stuff with last space
						text += line[0:(i - 1)]
						width += lineWidth
						height += lineHeight
						break
					}

					// If max for a line reached (and the text should wrap), make sure to go back
					if !node.Props().wrapping && (lineWidth > maxX || lineHeight > maxY) {
						// Cut line to stuff with last space
						line = line[0:(lastSpace - lineOffset)]
						text += line + "\n"
						line = ""
						lineOffset = lastSpace - lineOffset
						i = lineOffset
						continue
					}

					if unicode.IsSpace(rune) {
						lastSpace = i
					}

					i++
				}

				return scath.Vec{X: 0, Y: 0}, nil
			})

			props.Update(func(node *scaffui.SingleChildNode[TextProps], c *scaff.Context) (bool, error) {
				return false, nil
			})

			props.Draw(func(node *scaffui.SingleChildNode[TextProps], position scath.Vec, painter paint.Painter) {
				painter.Paint(paint.Text{
					Font:     node.Props().font,
					Text:     node.Props().text,
					Color:    node.Props().color,
					FontSize: node.Props().fontSize,
				})
			})
		},
	})
}

// This is extracted here so we can improve it in the future, with support for lots of different languages, etc.
func sliceIntoWords(text string) []string {
	return strings.Fields(text)
}
