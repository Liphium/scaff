package scaffui

import (
	"io/fs"

	"github.com/Liphium/scaff/paint"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type InterfaceLayer struct {
	root     *MountedNode
	assetsFS fs.FS
	renderer *paint.EbitenPainter
}

func NewInterfaceLayer(assetsFS fs.FS, builder NodeBuilder) *InterfaceLayer {
	return &InterfaceLayer{
		assetsFS: assetsFS,
		root:     NewMountedFromBuilder(builder),
	}
}

func (e *InterfaceLayer) Load() {
	e.root.Load(nil)
}

func (e *InterfaceLayer) Unload() {
	e.root.Unload()
}

func (e *InterfaceLayer) Draw(c *scaff.Context, screen *ebiten.Image) {
	// Initialize the UI and stuff
	if e.renderer == nil {
		e.renderer = paint.NewEbitenPainterWithFS(ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy()), true, e.assetsFS)
		e.root.Current().SetConstraints(Loose(c.Width, c.Height))
		_, err := e.root.Current().Layout()
		if err != nil {
			log.Error("layout error", "err", err)
		}
	}

	// Update all of the stuff
	result, err := e.root.Update(nil, c)
	if result.SizeChanged || err != nil {
		log.Warn("relayout or error happend", "result", result, "err", err)
		return
	}

	// Draw the stuff
	e.renderer.Clear()
	e.root.Current().Draw(scath.Zero, e.renderer)

	screen.DrawImage(e.renderer.Screen(), &ebiten.DrawImageOptions{})
}

func (e *InterfaceLayer) Update(c *scaff.Context) error {
	return nil
}
