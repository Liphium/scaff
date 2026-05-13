package basenode

import (
	"errors"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffui"
)

type ThrowerProps struct {
	msg string
}

func (tp *ThrowerProps) Message(message string) {
	tp.msg = message
}

func Thrower(create func(t *scaff.Tracker, tp *ThrowerProps)) scaffui.NodeBuilder {
	return scaffui.UseSingleNode("thrower", create, func(core *scaffui.SingleChildConstruct[ThrowerProps]) {
		core.Layout(func(node *scaffui.SingleChildNode[ThrowerProps]) (scaffui.Size, error) {
			return scaffui.Size{}, errors.New(node.Props().msg)
		})
	})
}
