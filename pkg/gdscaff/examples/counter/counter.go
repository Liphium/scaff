// This file is all you need to start a project.
// Save it somewhere, install the `gd` command and use `gd run` to launch.
package main

import (
	"fmt"

	"graphics.gd/startup"

	"graphics.gd/classdb/Button"
	"graphics.gd/classdb/CenterContainer"
	"graphics.gd/classdb/Control"
	"graphics.gd/classdb/SceneTree"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/pkg/gdscaff"
)

var counter = scaff.NewSignal(0)

func main() {
	// This would later be in some other package
	m := gdscaff.NewMount()

	gdscaff.With(m, CenterContainer.New).Init(func(t *scaff.Tracker, node CenterContainer.Instance) {

	})

	startup.LoadingScene() // setup the SceneTree and wait until we have access to engine functionality

	// TODO: This will also be enhanced by a dedicated package to make sure we have a nice tree structure
	center := CenterContainer.New()
	center.AsControl().SetAnchorsPreset(Control.PresetFullRect)

	button := Button.New()

	t := m.Instance.NewTracker(nil)

	button.AsBaseButton().OnPressed(func() {
		counter.Set(counter.Value() + 1)
	})

	t.Effect(func() {
		button.SetText(fmt.Sprintf("Count: %d", counter.Track(t)))
	})

	center.AsNode().AddChild(button.AsNode())
	m.AsNode().AddChild(center.AsNode())

	SceneTree.Add(m)
	startup.Scene() // starts up the scene and blocks until the engine shuts down.
}
