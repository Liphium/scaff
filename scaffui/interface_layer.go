package scaffui

import (
	"io/fs"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/smath"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var _ scaff.UILayer = &InterfaceLayer{}

type InterfaceLayer struct {
	root     *MountedNode
	assetsFS fs.FS
	renderer *EbitenRenderer

	prevCursorSet bool
	prevCursorX   int
	prevCursorY   int
}

func NewInterfaceLayer(assetsFS fs.FS, builder NodeBuilder) *InterfaceLayer {
	return &InterfaceLayer{
		assetsFS: assetsFS,
		root:     NewMountedFromBuilder(builder),
	}
}

func (e *InterfaceLayer) Load() {
	e.root.Load()
}

func (e *InterfaceLayer) Unload() {
	e.root.Unload()
}

func (e *InterfaceLayer) Draw(c *scaff.LayerContext, screen *ebiten.Image) {
	// Initialize the UI and stuff
	if e.renderer == nil {
		e.renderer = NewEbitenRendererWithFS(ebiten.NewImage(c.Width, c.Height), true, e.assetsFS)
		e.root.Current().SetConstraints(Loose(c.Width, c.Height))
		_, err := e.root.Current().Layout()
		if err != nil {
			log.Error("layout error", "err", err)
		}
	}

	x, y := ebiten.CursorPosition()

	deltaX := x - e.prevCursorX
	deltaY := y - e.prevCursorY
	e.prevCursorX = x
	e.prevCursorY = y

	if deltaX != 0 || deltaY != 0 {
		e.root.Current().HandleEvent(c, MoveEvent{
			X:      x,
			Y:      y,
			DeltaX: deltaX,
			DeltaY: deltaY,
		})
	}

	// Update all of the stuff
	relayout, err := e.root.Update(c)
	if relayout || err != nil {
		log.Warn("relayout or error happend", "relayout", relayout, "err", err)
		return
	}

	// Draw the stuff
	e.renderer.Clear()
	e.root.Current().Draw(smath.Zero, e.renderer)

	screen.DrawImage(e.renderer.Screen(), &ebiten.DrawImageOptions{})
}

func (e *InterfaceLayer) Update(c *scaff.LayerContext) error {
	if e.renderer != nil {
		x, y := ebiten.CursorPosition()

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			e.root.Current().HandleEvent(c, ReleaseEvent{
				X:      x,
				Y:      y,
				Button: LeftClick,
			})
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
			e.root.Current().HandleEvent(c, ReleaseEvent{
				X:      x,
				Y:      y,
				Button: RightClick,
			})
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
			e.root.Current().HandleEvent(c, ReleaseEvent{
				X:      x,
				Y:      y,
				Button: MiddleClick,
			})
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			e.root.Current().HandleEvent(c, DownEvent{
				X:      x,
				Y:      y,
				Button: LeftClick,
			})
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			e.root.Current().HandleEvent(c, DownEvent{
				X:      x,
				Y:      y,
				Button: RightClick,
			})
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
			e.root.Current().HandleEvent(c, DownEvent{
				X:      x,
				Y:      y,
				Button: MiddleClick,
			})
		}

		scrollX, scrollY := ebiten.Wheel()
		if scrollX != 0 || scrollY != 0 {
			e.root.Current().HandleEvent(c, ScrollEvent{
				X:       x,
				Y:       y,
				ScrollX: scrollX,
				ScrollY: scrollY,
			})
		}
	}
	return nil
}
