package scaff

import (
	"slices"
	"sync"
	"time"

	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/smath"
	"github.com/Liphium/scaff/util"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	mu               *sync.Mutex
	waitingSceneList []Scene
	modified         bool
	width, height    int

	sceneList []*TransitioningState[Scene] // THIS SHOULD NOT BE SET OUTSIDE OF UPDATE
}

// NewGame returns a new Game instance with default values
func NewGame() *Game {
	return &Game{
		mu:               &sync.Mutex{},
		waitingSceneList: []Scene{},
		sceneList:        []*TransitioningState[Scene]{},
	}
}

func (g *Game) Layout(width, height int) (int, int) {
	g.width = width
	g.height = height
	return width, height
}

func (g *Game) Update() error {
	now := time.Now()

	g.mu.Lock()

	// Remove scenes starting from the front of the scene list
	if g.modified {

		// Compare the lists and load new scenes / unload old oones
		for i, scene := range g.waitingSceneList {
			if i < len(g.sceneList) {

				// We specifically compare the scene pointers, not the scene IDs (to make sure that when a new instance is inserted in the same index, it still gets loaded)
				if val, ok := g.sceneList[i].GetCurrent().Value(); !ok || scene != val {
					g.sceneList[i].Set(now, optional.With(scene))
				}
			} else {
				g.sceneList = append(g.sceneList, NewTransitioningState(now, scene))
			}
		}

		// If the size of the sceneList is greater than the waitingSceneList, set excess scenes from the back to empty
		if len(g.sceneList) > len(g.waitingSceneList) {
			for i := len(g.waitingSceneList); i < len(g.sceneList); i++ {
				g.sceneList[i].SetEmpty(now)
			}
		}
	}

	g.mu.Unlock()

	// Update all of the scenes in proper order (this is backward because the front of the scene list is the topmost scene)
	for i, scene := range slices.Backward(g.sceneList) {
		if err := scene.Update(now, func(s Scene, tf smath.Timeframe) error {
			return s.Update(g.buildSceneContext(i, now, util.Ptr(tf)))
		}); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	now := time.Now()

	// Draw all of the scenes in proper order (this is forward because the front of the scene list is the scene that is in the background)
	for i, scene := range g.sceneList {
		scene.Update(now, func(s Scene, tf smath.Timeframe) error {
			s.Draw(g.buildSceneContext(i, now, util.Ptr(tf)), screen)
			return nil
		})
	}
}

func (g *Game) buildSceneContext(i int, now time.Time, frame *smath.Timeframe) SceneContext {

	// Create default frame if not set
	if frame == nil {
		// This is the default to make sure the transition is immediately over and no transition occurs
		frame = util.Ptr(smath.NewTimeframe(now, 0*time.Nanosecond))
	}

	return SceneContext{
		Focused:         i == len(g.sceneList)-1,
		Now:             now,
		TransitionFrame: *frame,
		Width:           g.width,
		Height:          g.height,
	}
}

// Get the scene list, this may not be the last updated one, but the one that will be set next time Update is called.
func (g *Game) GetSceneList() []Scene {
	g.mu.Lock()
	defer g.mu.Unlock()

	return slices.Clone(g.waitingSceneList)
}

func (g *Game) Goto(scene Scene) {
	g.mu.Lock()
	g.waitingSceneList = []Scene{scene}
	g.modified = true
	g.mu.Unlock()
}

func (g *Game) Push(scene Scene) {
	g.mu.Lock()
	g.waitingSceneList = append(g.waitingSceneList, scene)
	g.modified = true
	g.mu.Unlock()
}

func (g *Game) Pop() {
	g.mu.Lock()
	if len(g.waitingSceneList) > 0 {
		g.waitingSceneList = g.waitingSceneList[:len(g.waitingSceneList)-1]
		g.modified = true
	}
	g.mu.Unlock()
}

func (g *Game) PopUntil(id string) {
	g.mu.Lock()
	for len(g.waitingSceneList) > 0 {
		if g.waitingSceneList[len(g.waitingSceneList)-1].GetId() == id {
			break
		}
		g.waitingSceneList = g.waitingSceneList[:len(g.waitingSceneList)-1]
	}
	g.modified = true
	g.mu.Unlock()
}
