package gdscaff

import (
	"github.com/Liphium/scaff"
	"graphics.gd/classdb/Node"
)

type GodotNode interface {
	AsNode() Node.Instance
}

var _ scaff.Tracking = Ctx[Node.Instance]{}

type Ctx[N GodotNode] struct {
	mount   *Mount
	current *Reactor[N]
	tracker *scaff.Tracker
}

func (c Ctx[N]) Tracker() *scaff.Tracker {
	return c.tracker
}

type Reactor[N GodotNode] struct {
	create func() N
	init   func(t *Ctx[N], node N)

	mount   *Mount
	tracker *scaff.Tracker
}

func Child[P, C GodotNode](ctx Ctx[P], create func() C) *Reactor[C] {
	// TODO: Add build hook that adds child
	return With(ctx.mount, create)
}

func With[N GodotNode](m *Mount, create func() N) *Reactor[N] {
	return &Reactor[N]{
		create: create,

		mount:   m,
		tracker: m.Instance.NewTracker(nil),
	}
}

func (r *Reactor[N]) Init(init func(t *Ctx[N], node N)) {
	r.init = init
}

func (r *Reactor[N]) Build() N {
	if r.create == nil {
		panic("no creation function")
	}
	node := r.create()

	if r.init == nil {
		return node
	}

	r.init(r.tracker, node)
	node.AsNode().OnTreeExited(r.Unload)
	return node
}

func (r *Reactor[N]) buildCtx() *Ctx[N] {
	return &Ctx[N]{
		mount:   r.mount,
		current: r,
		tracker: r.tracker,
	}
}

func (r *Reactor[N]) Unload() {

}
