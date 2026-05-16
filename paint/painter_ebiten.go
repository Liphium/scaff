package paint

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var _ Painter = &EbitenPainter{}

// Create a new renderer using Ebitengine's image API (this uses the vector package)
func NewEbitenPainter(screen *ebiten.Image, antialias bool) *EbitenPainter {
	return NewEbitenPainterWithFS(screen, antialias, nil)
}

// Create a new renderer using Ebitengine's image API and a custom filesystem for assets.
func NewEbitenPainterWithFS(screen *ebiten.Image, antialias bool, assetsFS fs.FS) *EbitenPainter {
	if assetsFS == nil {
		assetsFS = os.DirFS(".")
	}

	return &EbitenPainter{
		screen:    screen,
		antialias: antialias,
		assetsFS:  assetsFS,
	}
}

type EbitenPainter struct {
	screen    *ebiten.Image
	antialias bool
	images    map[string]*ebiten.Image
	assetsFS  fs.FS
}

func (er *EbitenPainter) Screen() *ebiten.Image {
	return er.screen
}

func (er *EbitenPainter) Clear() {
	er.screen.Clear()
}

func (er *EbitenPainter) Paint(command RenderCommand) {
	if er.images == nil {
		er.images = make(map[string]*ebiten.Image)
	}

	switch c := command.(type) {
	case Rectangle:
		er.drawRectangle(c)
	case RectangleStroke:
		er.drawRectangleStroke(c)
	case Image:
		er.drawImage(c)
	case Text:
		log.Warn("text command not implemented in ebiten renderer yet")
	default:
		log.Warn("unknown render command", "id", command.ID())
	}
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

	path := roundedRectPath(command.Position.X, command.Position.Y, command.Size.X, command.Size.Y, command.BorderRadius)

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

	img := er.images[command.Path]
	if img == nil {
		loaded, err := er.loadEbitenImage(command.Path)
		if err != nil {
			log.Error("failed to load image", "path", command.Path, "err", err)
			return
		}
		img = loaded
		er.images[command.Path] = img
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
	er.screen.DrawImage(img, opts)
}

func roundedRectPath(x, y, width, height float64, borderRadius int) *vector.Path {
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

func (er *EbitenPainter) loadEbitenImage(path string) (*ebiten.Image, error) {
	f, err := er.assetsFS.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}
