package paint

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var _ Painter = &EbitenPainter{}

// Create a new renderer using Ebitengine's image API and a custom filesystem for assets.
func NewEbitenPainter(screen *ebiten.Image, antialias bool, assets *AssetManager) *EbitenPainter {
	return &EbitenPainter{
		screen:    screen,
		antialias: antialias,
		assets:    assets,
		transform: NewTransform(),
	}
}

type EbitenPainter struct {
	transform Transform
	screen    *ebiten.Image
	antialias bool
	assets    *AssetManager
}

func (e *EbitenPainter) SetTransform(transform Transform) {
	e.transform = transform
}

func (e *EbitenPainter) Transform() Transform {
	return e.transform
}

func (er *EbitenPainter) Screen() *ebiten.Image {
	return er.screen
}

func (er *EbitenPainter) Clear() {
	er.screen.Clear()
}

func (er *EbitenPainter) Paint(command RenderCommand) {
	switch c := command.(type) {
	case Rectangle:
		er.drawRectangle(c)
	case RectangleStroke:
		er.drawRectangleStroke(c)
	case Image:
		er.drawImage(c)
	case Text:
		er.drawText(c)
	default:
		log.Warn("unknown render command", "id", command.ID())
	}
}

func (er *EbitenPainter) DrawRaw(image *ebiten.Image, op *ebiten.DrawImageOptions) {

	// Apply painter transform
	op.GeoM.Translate(-er.transform.CamX, -er.transform.CamY)
	op.GeoM.Translate(er.transform.CenterOffsetX, er.transform.CenterOffsetY)
	op.GeoM.Rotate(er.transform.Angle)
	op.GeoM.Scale(er.transform.ZoomFactor, er.transform.ZoomFactor)
	op.GeoM.Translate(math.Abs(er.transform.CenterOffsetX), math.Abs(er.transform.CenterOffsetY))

	er.screen.DrawImage(image, op)
}

func (er *EbitenPainter) PaintMulti(commands []RenderCommand) {
	for _, command := range commands {
		er.Paint(command)
	}
}

func (er *EbitenPainter) drawRectangle(command Rectangle) {
	if command.Size.X <= 0 || command.Size.Y <= 0 {
		return
	}

	path := roundedRectPath(command.Position.X, command.Position.Y, command.Size.X, command.Size.Y, command.BorderRadius)

	// Apply transform to path
	if er.transform.ZoomFactor != 1 || er.transform.Angle != 0 || er.transform.CamX != 0 || er.transform.CamY != 0 || er.transform.CenterOffsetX != 0 || er.transform.CenterOffsetY != 0 {
		geom := ebiten.GeoM{}
		geom.Translate(-er.transform.CamX, -er.transform.CamY)
		geom.Translate(er.transform.CenterOffsetX, er.transform.CenterOffsetY)
		geom.Rotate(er.transform.Angle)
		geom.Scale(er.transform.ZoomFactor, er.transform.ZoomFactor)
		geom.Translate(math.Abs(er.transform.CenterOffsetX), math.Abs(er.transform.CenterOffsetY))

		tfPath := &vector.Path{}
		tfPath.AddPath(path, &vector.AddPathOptions{GeoM: geom})
		path = tfPath
	}

	drawOptions := &vector.DrawPathOptions{AntiAlias: er.antialias}
	drawOptions.ColorScale.ScaleWithColor(command.FillColor)
	vector.FillPath(er.screen, path, nil, drawOptions)
}

func (er *EbitenPainter) drawRectangleStroke(command RectangleStroke) {
	if command.Size.X <= 0 || command.Size.Y <= 0 {
		return
	}

	thickness := command.Thickness
	if thickness <= 0 {
		thickness = 1
	}

	path := roundedRectPath(command.Position.X+command.Thickness/2, command.Position.Y+command.Thickness/2, command.Size.X-command.Thickness, command.Size.Y-command.Thickness, command.BorderRadius)

	// Apply transform to path
	if er.transform.ZoomFactor != 1 || er.transform.Angle != 0 || er.transform.CamX != 0 || er.transform.CamY != 0 || er.transform.CenterOffsetX != 0 || er.transform.CenterOffsetY != 0 {
		geom := ebiten.GeoM{}
		geom.Translate(-er.transform.CamX, -er.transform.CamY)
		geom.Translate(er.transform.CenterOffsetX, er.transform.CenterOffsetY)
		geom.Rotate(er.transform.Angle)
		geom.Scale(er.transform.ZoomFactor, er.transform.ZoomFactor)
		geom.Translate(math.Abs(er.transform.CenterOffsetX), math.Abs(er.transform.CenterOffsetY))

		tfPath := &vector.Path{}
		tfPath.AddPath(path, &vector.AddPathOptions{GeoM: geom})
		path = tfPath
	}

	strokeOptions := &vector.StrokeOptions{
		Width:    float32(thickness),
		LineJoin: vector.LineJoinRound,
		LineCap:  vector.LineCapRound,
	}
	drawOptions := &vector.DrawPathOptions{AntiAlias: er.antialias}
	drawOptions.ColorScale.ScaleWithColor(command.Color)
	vector.StrokePath(er.screen, path, strokeOptions, drawOptions)
}

func (er *EbitenPainter) drawImage(command Image) {
	if command.Size.X <= 0 || command.Size.Y <= 0 {
		return
	}

	img, err := er.assets.GetImage(command.Path)
	if err != nil {
		log.Error("failed to load image", "path", command.Path, "err", err)
		return
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w <= 0 || h <= 0 {
		return
	}

	opts := &ebiten.DrawImageOptions{}
	opts.Filter = command.FilterMode
	opts.GeoM.Scale(command.Size.X/float64(w), command.Size.Y/float64(h))
	opts.GeoM.Translate(command.Position.X, command.Position.Y)

	// Apply painter transform
	opts.GeoM.Translate(-er.transform.CamX, -er.transform.CamY)
	opts.GeoM.Translate(er.transform.CenterOffsetX, er.transform.CenterOffsetY)
	opts.GeoM.Rotate(er.transform.Angle)
	opts.GeoM.Scale(er.transform.ZoomFactor, er.transform.ZoomFactor)
	opts.GeoM.Translate(math.Abs(er.transform.CenterOffsetX), math.Abs(er.transform.CenterOffsetY))

	er.screen.DrawImage(img, opts)
}

func (er *EbitenPainter) drawText(command Text) {
	font := defaultFontFace

	if wanted, err := er.assets.GetFont(command.Font); err == nil {
		font = wanted
	} else {
		log.Warn("font not found, using default font", "name", command.Font, "text", command.Text)
	}

	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			LineSpacing:    (command.FontSize + command.LineSpacing) * er.transform.ZoomFactor,
			PrimaryAlign:   command.PrimaryAlign,
			SecondaryAlign: command.SecondaryAlign,
		},
	}
	op.GeoM.Translate(command.Position.X, command.Position.Y)

	// Apply rotation and translation from transform (scaling is done through font size to keep text sharp)
	op.GeoM.Translate(-er.transform.CamX, -er.transform.CamY)
	op.GeoM.Translate(er.transform.CenterOffsetX, er.transform.CenterOffsetY)
	op.GeoM.Rotate(er.transform.Angle)
	// We do NOT scale the geom for text to prevent blurriness.
	op.GeoM.Translate(math.Abs(er.transform.CenterOffsetX), math.Abs(er.transform.CenterOffsetY))

	op.ColorScale.ScaleWithColor(command.Color)

	// Multiply font size by ZoomFactor to render sharply at correct scale
	fontSize := command.FontSize * er.transform.ZoomFactor
	if fontSize <= 0 {
		return
	}

	text.Draw(er.screen, command.Text, &text.GoTextFace{
		Source:    font,
		Direction: command.Direction,
		Size:      fontSize,
	}, op)
}

func roundedRectPath(x, y, width, height float64, borderRadius float64) *vector.Path {
	path := &vector.Path{}
	if width <= 0 || height <= 0 {
		return path
	}

	r := float32(borderRadius)
	if r <= 0 {
		path.MoveTo(float32(x), float32(y))
		path.LineTo(float32(x+width), float32(y))
		path.LineTo(float32(x+width), float32(y+height))
		path.LineTo(float32(x), float32(y+height))
		path.Close()
		return path
	}

	maxRadius := float32(math.Min(width, height) / 2)
	if r > maxRadius {
		r = maxRadius
	}

	left := float32(x)
	top := float32(y)
	right := float32(x + width)
	bottom := float32(y + height)

	path.MoveTo(left+r, top)
	path.LineTo(right-r, top)
	path.Arc(right-r, top+r, r, -float32(math.Pi)/2, 0, vector.Clockwise)
	path.LineTo(right, bottom-r)
	path.Arc(right-r, bottom-r, r, 0, float32(math.Pi)/2, vector.Clockwise)
	path.LineTo(left+r, bottom)
	path.Arc(left+r, bottom-r, r, float32(math.Pi)/2, float32(math.Pi), vector.Clockwise)
	path.LineTo(left, top+r)
	path.Arc(left+r, top+r, r, float32(math.Pi), float32(math.Pi)*3/2, vector.Clockwise)
	path.Close()

	return path
}
