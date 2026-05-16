package basenode

import (
	"errors"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scath"
)

type ThrowerProps struct {
	msg string
}

func (tp *ThrowerProps) Message(message string) {
	tp.msg = message
}

func Thrower(create func(t *scaff.Tracker, tp *ThrowerProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode("thrower", create, func(core *scaffui.SingleChildProps[ThrowerProps]) {
		core.Layout(func(node *scaffui.SingleChildNode[ThrowerProps]) (scath.Vec, error) {
			return scath.Vec{}, errors.New(node.Props().msg)
		})
	})
}
