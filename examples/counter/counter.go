// This file is all you need to start a project.
// Save it somewhere, install the `gd` command and use `gd run` to launch.
package main

import (
	"log"

	"graphics.gd/classdb"
	"graphics.gd/startup"
	"graphics.gd/variant/Float"

	"graphics.gd/classdb/Control"
	"graphics.gd/classdb/GUI"
	"graphics.gd/classdb/Label"
	"graphics.gd/classdb/Node"
	"graphics.gd/classdb/SceneTree"
)

type MyCustomNode struct {
	Node.Extension[MyCustomNode]
}

func (m *MyCustomNode) Process(delta Float.X) {
	log.Println("hello world")
}

func main() {
	startup.LoadingScene() // setup the SceneTree and wait until we have access to engine functionality
	hello := Label.New()
	hello.AsControl().SetAnchorsPreset(Control.PresetFullRect) // expand the label to take up the whole screen.
	hello.SetHorizontalAlignment(GUI.HorizontalAlignmentCenter)
	hello.SetVerticalAlignment(GUI.VerticalAlignmentCenter)
	hello.SetText("Hello, World!")

	classdb.Register[MyCustomNode]()
	hello.AsNode().AddChild(new(MyCustomNode).AsNode())
	SceneTree.Add(hello)
	startup.Scene() // starts up the scene and blocks until the engine shuts down.
}
