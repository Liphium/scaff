package gdscaff

import (
	"sync"
	"sync/atomic"

	"github.com/Liphium/scaff"
	"graphics.gd/classdb/Node"
)

var counter = atomic.Int32{}

// id (int32) -> *scaff.Instance
var instanceMap = sync.Map{}

// A node for mounting scaff with an instance. This is the preferred way to interact with Scaff when using Godot.
type Mount struct {
	id int32
	Node.Extension[Mount]
}

// Create a new scaff instance when the node is intialized
func (m *Mount) Ready() {
	m.id = counter.Add(1)
	instanceMap.Store(m.id, scaff.NewInstance())
}

// When the node is deleted, make sure to delete the instance as well
func (m *Mount) OnTreeExit() {
	Instance(m).Unload()
	instanceMap.Delete(m.id)
}

// In the regular update loop, run all trackers.
func (m *Mount) Process() {
	Instance(m).Update()
}

// Get the scaff.Instance for a mount.
func Instance(m *Mount) *scaff.Instance {
	if v, ok := instanceMap.Load(m.id); ok {
		return v.(*scaff.Instance)
	}
	panic("instance not found for mount, not initialized yet?")
}
