package gdscaff

import (
	"github.com/Liphium/scaff"
	sutil "github.com/Liphium/scaff/util"
	"graphics.gd/classdb/Node"
)

var log = sutil.NewLogger("gdscaff")

type Context[N any] struct{}

type GodotNode interface {
	AsNode() Node.Instance
}

func With[N GodotNode](creator func() N) *Context[N] {
	return &Context[N]{}
}

func (c *Context[N]) Build(create func(t *scaff.Tracker, node N)) {

}
