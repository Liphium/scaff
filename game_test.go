package scaff_test

import (
	"testing"

	"github.com/Liphium/scaff"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

var _ scaff.Scene = &TestScene{}

type TestScene struct {
	ID             string
	OnUpdate       func(c *scaff.Context) error
	OnLoad         func()
	OnUnload       func()
	TransitionFunc func(in bool) scaff.TransitionProperties
}

func (s *TestScene) GetId() string {
	return s.ID
}

func (s *TestScene) Update(c *scaff.Context) error {
	return s.OnUpdate(c)
}

func (s *TestScene) Draw(c *scaff.Context, screen *ebiten.Image) {}

func (s *TestScene) Load() {
	if s.OnLoad != nil {
		s.OnLoad()
	}
}
func (s *TestScene) Unload() {
	if s.OnUnload != nil {
		s.OnUnload()
	}
}

func (s *TestScene) Transition(in bool) scaff.TransitionProperties {
	if s.TransitionFunc != nil {
		return s.TransitionFunc(in)
	}
	return scaff.NoTransition()
}

func createMockScene(id string, updates *[]string, loads *[]string, unloads *[]string) *TestScene {
	return &TestScene{
		ID: id,
		OnUpdate: func(c *scaff.Context) error {
			if updates != nil {
				*updates = append(*updates, id)
			}
			return nil
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

// Tests that pushing a scene onto the game works as expected.
func TestGamePush(t *testing.T) {
	g := scaff.NewGame()

	g.Push(createMockScene("test1", nil, nil, nil))
	g.Update()

	assert.Equal(t, 1, len(g.GetSceneList()))
	assert.Equal(t, "test1", g.GetSceneList()[0].GetId())
}

// Tests that popping a scene from the game works as expected.
func TestGamePop(t *testing.T) {
	g := scaff.NewGame()

	g.Push(createMockScene("test1", nil, nil, nil))
	g.Pop()
	g.Update()

	assert.Equal(t, 0, len(g.GetSceneList()))
}

// Tests that popping all scenes until a specific scene is popped works as expected.
func TestGamePopUntil(t *testing.T) {
	g := scaff.NewGame()

	g.Push(createMockScene("test1", nil, nil, nil))
	g.Push(createMockScene("test2", nil, nil, nil))
	g.Push(createMockScene("test3", nil, nil, nil))
	g.Push(createMockScene("test4", nil, nil, nil))
	g.PopUntil("test3")
	g.Update()

	assert.Equal(t, 3, len(g.GetSceneList()))
	assert.Equal(t, "test1", g.GetSceneList()[0].GetId())
	assert.Equal(t, "test2", g.GetSceneList()[1].GetId())
	assert.Equal(t, "test3", g.GetSceneList()[2].GetId())
}

// Tests that the goto method works as expected.
func TestGameGoto(t *testing.T) {
	g := scaff.NewGame()

	g.Push(createMockScene("test1", nil, nil, nil))
	g.Push(createMockScene("test2", nil, nil, nil))
	g.Update()
	g.Goto(createMockScene("test3", nil, nil, nil))
	g.Update()

	assert.Equal(t, 1, len(g.GetSceneList()))
	assert.Equal(t, "test3", g.GetSceneList()[0].GetId())
}

// Makes sure scenes are loaded, unloaded and updated in the correct order.
func TestGameUpdate(t *testing.T) {
	t.Run("scenes pushed are in correct order", func(t *testing.T) {
		g := scaff.NewGame()
		updates := []string{}
		loads := []string{}
		unloads := []string{}

		g.Push(createMockScene("scene1", &updates, &loads, &unloads))
		g.Push(createMockScene("scene2", &updates, &loads, &unloads))
		g.Update()

		assert.Equal(t, []string{"scene2", "scene1"}, updates)
		assert.Equal(t, []string{"scene1", "scene2"}, loads)
		assert.Equal(t, []string{}, unloads)
	})

	t.Run("scenes updated in correct order", func(t *testing.T) {
		g := scaff.NewGame()
		updates := []string{}
		loads := []string{}
		unloads := []string{}

		g.Push(createMockScene("scene1", &updates, &loads, &unloads))
		g.Update()
		g.Push(createMockScene("scene2", &updates, &loads, &unloads))
		g.Update()

		assert.Equal(t, []string{"scene1", "scene2", "scene1"}, updates)
		assert.Equal(t, []string{"scene1", "scene2"}, loads)
		assert.Equal(t, []string{}, unloads)
	})

	t.Run("scenes are properly unloaded", func(t *testing.T) {
		g := scaff.NewGame()
		updates := make([]string, 0, 3)
		loads := make([]string, 0, 2)
		unloads := make([]string, 0, 1)

		g.Push(createMockScene("scene1", &updates, &loads, &unloads))
		g.Push(createMockScene("scene2", &updates, &loads, &unloads))
		g.Update()
		g.Pop()
		g.Update()

		assert.Equal(t, []string{"scene2", "scene1", "scene1"}, updates)
		assert.Equal(t, []string{"scene1", "scene2"}, loads)
		assert.Equal(t, []string{"scene2"}, unloads)
	})

	t.Run("scenes can be replaced", func(t *testing.T) {
		g := scaff.NewGame()
		updates := []string{}
		loads := []string{}
		unloads := []string{}

		scene1 := createMockScene("scene1", &updates, &loads, &unloads)
		scene2 := createMockScene("scene2", &updates, &loads, &unloads)
		scene3 := createMockScene("scene3", &updates, &loads, &unloads)

		g.Push(scene1)
		g.Push(scene2)
		g.Update()
		g.Pop()
		g.Push(scene3)
		g.Update()

		assert.Equal(t, []string{"scene2", "scene1", "scene3", "scene1"}, updates)
		assert.Equal(t, []string{"scene1", "scene2", "scene3"}, loads)
		assert.Equal(t, []string{"scene2"}, unloads)
	})

	t.Run("scenes can be replaced, 2 layers", func(t *testing.T) {
		g := scaff.NewGame()
		updates := []string{}
		loads := []string{}
		unloads := []string{}

		scene1 := createMockScene("scene1", &updates, &loads, &unloads)
		scene2 := createMockScene("scene2", &updates, &loads, &unloads)
		scene3 := createMockScene("scene3", &updates, &loads, &unloads)
		scene4 := createMockScene("scene4", &updates, &loads, &unloads)

		g.Push(scene1)
		g.Push(scene2)
		g.Push(scene3)
		g.Update()
		g.PopUntil("scene1")
		g.Push(scene4)
		g.Push(scene3)
		g.Update()

		assert.Equal(t, []string{"scene3", "scene2", "scene1", "scene3", "scene4", "scene1"}, updates)
		assert.Equal(t, []string{"scene1", "scene2", "scene3", "scene4"}, loads)
		assert.Equal(t, []string{"scene2"}, unloads)
	})

	t.Run("replacement with scene of same id still fires load", func(t *testing.T) {
		g := scaff.NewGame()
		updates := []string{}
		loads := []string{}
		unloads := []string{}

		scene1 := createMockScene("scene1", &updates, &loads, &unloads)
		scene2 := createMockScene("scene2", &updates, &loads, &unloads)
		scene22 := createMockScene("scene2", &updates, &loads, &unloads)

		g.Push(scene1)
		g.Push(scene2)
		g.Update()
		g.Pop()
		g.Push(scene22)
		g.Update()

		assert.Equal(t, []string{"scene2", "scene1", "scene2", "scene1"}, updates)
		assert.Equal(t, []string{"scene1", "scene2", "scene2"}, loads)
		assert.Equal(t, []string{"scene2"}, unloads)
	})
}
