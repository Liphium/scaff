package scaff_test

import (
	"testing"

	"github.com/Liphium/scaff"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

var _ scaff.WorldLayer = &TestWorldLayer{}
var _ scaff.UILayer = &TestUILayer{}

type TestWorldLayer struct {
	ID       string
	OnUpdate func(c *scaff.LayerContext) error
	OnDraw   func(c *scaff.LayerContext, camera *scaff.Camera, screen *ebiten.Image)
	OnLoad   func()
	OnUnload func()
}

func (l *TestWorldLayer) Update(c *scaff.LayerContext) error {
	if l.OnUpdate != nil {
		return l.OnUpdate(c)
	}
	return nil
}

func (l *TestWorldLayer) Draw(c *scaff.LayerContext, camera *scaff.Camera, screen *ebiten.Image) {
	if l.OnDraw != nil {
		l.OnDraw(c, camera, screen)
	}
}

func (l *TestWorldLayer) Load() {
	if l.OnLoad != nil {
		l.OnLoad()
	}
}

func (l *TestWorldLayer) Unload() {
	if l.OnUnload != nil {
		l.OnUnload()
	}
}

type TestUILayer struct {
	ID       string
	OnUpdate func(c *scaff.LayerContext) error
	OnDraw   func(c *scaff.LayerContext, screen *ebiten.Image)
	OnLoad   func()
	OnUnload func()
}

func (l *TestUILayer) Update(c *scaff.LayerContext) error {
	if l.OnUpdate != nil {
		return l.OnUpdate(c)
	}
	return nil
}

func (l *TestUILayer) Draw(c *scaff.LayerContext, screen *ebiten.Image) {
	if l.OnDraw != nil {
		l.OnDraw(c, screen)
	}
}

func (l *TestUILayer) Load() {
	if l.OnLoad != nil {
		l.OnLoad()
	}
}

func (l *TestUILayer) Unload() {
	if l.OnUnload != nil {
		l.OnUnload()
	}
}

func createMockWorldLayer(id string, updates *[]string, draws *[]string, loads *[]string, unloads *[]string) *TestWorldLayer {
	return &TestWorldLayer{
		ID: id,
		OnUpdate: func(c *scaff.LayerContext) error {
			if updates != nil {
				*updates = append(*updates, id)
			}
			return nil
		},
		OnDraw: func(c *scaff.LayerContext, camera *scaff.Camera, screen *ebiten.Image) {
			if draws != nil {
				*draws = append(*draws, id)
			}
		},
		OnLoad: func() {
			if loads != nil {
				*loads = append(*loads, id)
			}
		},
		OnUnload: func() {
			if unloads != nil {
				*unloads = append(*unloads, id)
			}
		},
	}
}

func createMockUILayer(id string, updates *[]string, draws *[]string, loads *[]string, unloads *[]string) *TestUILayer {
	return &TestUILayer{
		ID: id,
		OnUpdate: func(c *scaff.LayerContext) error {
			if updates != nil {
				*updates = append(*updates, id)
			}
			return nil
		},
		OnDraw: func(c *scaff.LayerContext, screen *ebiten.Image) {
			if draws != nil {
				*draws = append(*draws, id)
			}
		},
		OnLoad: func() {
			if loads != nil {
				*loads = append(*loads, id)
			}
		},
		OnUnload: func() {
			if unloads != nil {
				*unloads = append(*unloads, id)
			}
		},
	}
}

func TestOrder(t *testing.T) {
	updates, draws, loads, unloads := []string{}, []string{}, []string{}, []string{}

	scene := scaff.LayeredScene{
		ID: "test1",
		WorldLayers: []scaff.WorldLayer{
			createMockWorldLayer("wl1", &updates, &draws, &loads, &unloads),
			createMockWorldLayer("wl2", &updates, &draws, &loads, &unloads),
		},
		UILayers: []scaff.UILayer{
			createMockUILayer("ul1", &updates, &draws, &loads, &unloads),
			createMockUILayer("ul2", &updates, &draws, &loads, &unloads),
		},
	}

	scene.Load()
	err := scene.Update(scaff.SceneContext{})
	assert.NoError(t, err)
	scene.Draw(scaff.SceneContext{Width: 1920, Height: 1080}, nil)
	scene.Unload()

	assert.Equal(t, []string{"wl1", "wl2", "ul1", "ul2"}, loads)
	assert.Equal(t, []string{"ul2", "ul1", "wl2", "wl1"}, updates)
	assert.Equal(t, []string{"wl1", "wl2", "ul1", "ul2"}, draws)
	assert.Equal(t, []string{"wl1", "wl2", "ul1", "ul2"}, unloads)
}

func TestEventPropagation(t *testing.T) {
	scene := scaff.LayeredScene{
		ID: "test-propagation",
		WorldLayers: []scaff.WorldLayer{
			&TestWorldLayer{
				ID: "world-bottom",
				OnUpdate: func(c *scaff.LayerContext) error {
					assert.True(t, c.IsHandled(scaff.EventMouseInput))
					return nil
				},
			},
		},
		UILayers: []scaff.UILayer{
			&TestUILayer{
				ID: "ui-top",
				OnUpdate: func(c *scaff.LayerContext) error {
					c.Handled(scaff.EventMouseInput)
					return nil
				},
			},
		},
	}

	err := scene.Update(scaff.SceneContext{})
	assert.NoError(t, err)
}
