// This file is all you need to start a project.
// Save it somewhere, install the `gd` command and use `gd run` to launch.
package main

import (
	"fmt"

	"graphics.gd/classdb"
	"graphics.gd/startup"
	"graphics.gd/variant/Float"

	"graphics.gd/classdb/Button"
	"graphics.gd/classdb/CenterContainer"
	"graphics.gd/classdb/Control"
	"graphics.gd/classdb/SceneTree"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/pkg/gdscaff"
)

var Instance = scaff.NewInstance()

type MyCustomNode struct {
	Control.Extension[MyCustomNode]
}

func (m *MyCustomNode) Process(delta Float.X) {
	Instance.Update()
}

var counter = scaff.NewSignal(0)

func main() {
	// This would later be in some other package
	classdb.Register[MyCustomNode]()
	node := new(MyCustomNode)
	t := Instance.NewTracker(nil) // Trackers would be per node, but for now this works too

	startup.LoadingScene() // setup the SceneTree and wait until we have access to engine functionality

	// TODO: This will also be enhanced by a dedicated package to make sure we have a nice tree structure
	center := CenterContainer.New()
	center.AsControl().SetAnchorsPreset(Control.PresetFullRect)

	button := Button.New()
	t.Effect(func() {
		button.SetText(fmt.Sprintf("Count: %d", counter.Track(t)))
	})

	gdscaff.With(Button.New).Build(func(t *scaff.Tracker, node Button.Instance) {

	})

	button.AsBaseButton().OnPressed(func() {
		counter.Set(counter.Value() + 1)
	})

	center.AsNode().AddChild(button.AsNode())

	SceneTree.Add(center)
	SceneTree.Add(node)
	startup.Scene() // starts up the scene and blocks until the engine shuts down.
}
