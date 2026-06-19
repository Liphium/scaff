package gdscaff

import (
	"github.com/Liphium/scaff"
	"graphics.gd/classdb"
	"graphics.gd/classdb/Node"
)

// A node for mounting scaff with an instance. This is the preferred way to interact with Scaff when using Godot.
type Mount struct {
	Node.Extension[Mount]

	Instance *scaff.Instance `gd:"-"`
}

func init() {
	classdb.Register[Mount]()
}

func NewMount() *Mount {
	m := new(Mount)
	m.Instance = scaff.NewInstance()
	return m
}

// When the node is deleted, make sure to delete the instance as well
func (m *Mount) OnTreeExit() {
	m.Instance.Unload()
}

// In the regular update loop, run all trackers.
func (m *Mount) Process(delta float32) {
	m.Instance.Update()
}
