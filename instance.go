package scaff

type Instance struct {
	updateQueue *UpdateQueue
}

func NewInstance() *Instance {
	return &Instance{
		updateQueue: NewUpdateQueue(),
	}
}

// Create a new tracker.
//
// You can pass in a new function that will run whenever something in the tracker has changed.
func (i *Instance) NewTracker(onChange func()) *Tracker {
	return newTracker(i, onChange)
}

func (i *Instance) Update() {
	i.updateQueue.Update()
}

func (i *Instance) Unload() {
	i.updateQueue.Clear()
}
