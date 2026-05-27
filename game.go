package scaff

import (
	"slices"
	"sync"
	"time"

	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scath"
	sutil "github.com/Liphium/scaff/util"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	mu               *sync.Mutex
	waitingSceneList []Scene
	modified         bool
	width, height    float64
	layoutChange     bool // To emit a layout change event in the Update when something changed

	sceneList []*TransitioningState[Scene] // THIS SHOULD NOT BE SET OUTSIDE OF UPDATE

	input *inputHandler
}

// NewGame returns a new Game instance with default values
func NewGame() *Game {
	return &Game{
		mu:               &sync.Mutex{},
		waitingSceneList: []Scene{},
		sceneList:        []*TransitioningState[Scene]{},
		input:            &inputHandler{},
	}
}

// This is not needed as we have LayoutF below, we want to make sure the scale factor is taken into account properly.
func (g *Game) Layout(width, height int) (int, int) {
	return -1, -1
}

// This makes sure the monitor's resolution is actually properly respected (Source: https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/ebitengine_game.md#layout)
func (g *Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.Monitor().DeviceScaleFactor()
	newWidth := logicWinWidth * scale
	newHeight := logicWinHeight * scale
	if newWidth != g.width || newHeight != g.height {
		g.layoutChange = true
	}
	g.width = newWidth
	g.height = newHeight
	return g.width, g.height
}

// Update forwards Ebiten's Update to all scenes and emits input events from top -> bottom.
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

	// Build all of the scene contexts for updating
	ctx := &eventContext{}
	contexts := g.buildSceneContexts(now, ctx)

	// Emit all of the events for the update pass
	for _, event := range g.input.getInputEventsUpdate() {
		g.EmitEvent(contexts, event)
	}

	// Emit layout change when one happened
	if g.layoutChange {
		g.EmitEvent(contexts, SizeChangeEvent{})
		g.layoutChange = false
	}

	// Update all of the scenes in proper order (this is backward because the front of the scene list is the topmost scene)
	for i, scene := range slices.Backward(g.sceneList) {
		if err := scene.Update(now, func(s Scene, tf scath.Timeframe) error {
			// This still requires one more build of the scene thingy since there may be multiple occurences of the same scene updating due to transitioning out and stuff
			return s.Update(g.buildSceneContext(i, now, sutil.Ptr(tf), ctx))
		}); err != nil {
			return err
		}
	}
	return nil
}

// Draw handles drawing of the game by drawing the bottom scene first and then going up from there.
func (g *Game) Draw(screen *ebiten.Image) {
	now := time.Now()

	// Build all of the scene contexts for drawing
	ctx := &eventContext{}
	contexts := g.buildSceneContexts(now, ctx)

	// Emit all of the events for the draw pass
	for _, event := range g.input.getInputEventsDraw() {
		g.EmitEvent(contexts, event)
	}

	// Draw all of the scenes in proper order (this is forward because the front of the scene list is the scene that is in the background)
	for i, scene := range g.sceneList {
		scene.Update(now, func(s Scene, tf scath.Timeframe) error {
			// This still requires one more build of the scene thingy since there may be multiple occurences of the same scene rendering due to transitioning out and stuff
			s.Draw(g.buildSceneContext(i, now, sutil.Ptr(tf), ctx), screen)
			return nil
		})
	}
}

// buildSceneContext builds a context for one scene based on its index in sceneList and other parameters from its transition.
func (g *Game) buildSceneContext(i int, now time.Time, frame *scath.Timeframe, ctx *eventContext) *Context {

	// Create default frame if not set
	if frame == nil {
		// This is the default to make sure the transition is immediately over and no transition occurs
		frame = sutil.Ptr(scath.NewTimeframe(now, 0*time.Nanosecond))
	}

	return &Context{
		focused:         i == len(g.sceneList)-1,
		now:             now,
		transitionFrame: *frame,
		width:           g.width,
		height:          g.height,
		eventContext:    ctx,
	}
}

// buildSceneContexts builds the context needed for updates, etc. for all scenes in the sceneList, the index matches the one in sceneList.
//
// There may be nil contexts for scenes that are currently transitioning out / are not there.
func (g *Game) buildSceneContexts(now time.Time, ctx *eventContext) []*Context {
	contexts := make([]*Context, len(g.sceneList))

	// Go through all scenes and build their context
	for i, scene := range g.sceneList {
		scene.Update(now, func(s Scene, tf scath.Timeframe) error {
			// We only want the currently active scene / transitioning in scene
			if tf.IsBackwards() {
				return nil
			}

			contexts[i] = g.buildSceneContext(i, now, sutil.Ptr(tf), ctx)
			return nil
		})
	}

	return contexts
}

// EmitEvent sends an event to all scenes in the scene list in proper order (top -> bottom)
func (g *Game) EmitEvent(contexts []*Context, event Event) {

	// Emit events in proper order (this is backward because the front of the scene list is the topmost scene)
	for i, scene := range slices.Backward(g.sceneList) {
		if val, ok := scene.GetCurrent().Value(); ok && contexts[i] != nil {
			if err := val.HandleEvent(contexts[i], event); err != nil {
				log.Error("error during input handling", "scene", val.GetId(), "err", err)
				return
			}
		}
	}
}

// EmitEventNoContext sends an event to all scenes in the scene list in proper order (top -> bottom), but it also builds a context for the scenes at the same time. Scenes that are not fully transitioned will not receive events.
func (g *Game) EmitEventNoContext(event Event) {
	ctx := &eventContext{}
	contexts := g.buildSceneContexts(time.Now(), ctx)
	g.EmitEvent(contexts, event)
}

// Get the scene list, this may not be the last updated one, but the one that will be set next time Update is called.
func (g *Game) GetSceneList() []Scene {
	g.mu.Lock()
	defer g.mu.Unlock()

	return slices.Clone(g.waitingSceneList)
}

// Goto completely removes all scenes from the scene stack and makes this scene the new one and only scene.
func (g *Game) Goto(scene Scene) {
	g.mu.Lock()
	g.waitingSceneList = []Scene{scene}
	g.modified = true
	g.mu.Unlock()
}

// Push adds a scene to the top of the scene stack.
func (g *Game) Push(scene Scene) {
	g.mu.Lock()
	g.waitingSceneList = append(g.waitingSceneList, scene)
	g.modified = true
	g.mu.Unlock()
}

// Pop removes the top scene from the scene stack.
func (g *Game) Pop() {
	g.mu.Lock()
	if len(g.waitingSceneList) > 0 {
		g.waitingSceneList = g.waitingSceneList[:len(g.waitingSceneList)-1]
		g.modified = true
	}
	g.mu.Unlock()
}

// PopUntil "pops" (removes from the top) all scenes in the scene stack until a scene with a certain id is found (but that one is not removed).
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
