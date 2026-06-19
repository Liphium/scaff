package gdscaff

import (
	"github.com/Liphium/scaff"
	sutil "github.com/Liphium/scaff/util"
	"graphics.gd/classdb/Node"
)

var log = sutil.NewLogger("gdscaff")

type GodotNode interface {
	AsNode() Node.Instance
}

func New[N GodotNode](creator func() N, create func(t *scaff.Tracker, node N)) {

}
