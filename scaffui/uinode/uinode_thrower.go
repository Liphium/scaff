package uinode

import (
	"errors"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scath"
)

type ThrowerProps struct {
	Msg string
	*scaffui.AcceptNoChild
}

func Thrower(create func(t *scaff.Tracker, tp *ThrowerProps)) scaffui.NodeBuilder {
	return scaffui.SingleNode(scaffui.SingleNodeCreate[ThrowerProps]{
		ID: "thrower",
		DefaultProps: ThrowerProps{
			AcceptNoChild: &scaffui.AcceptNoChild{},
			Msg:           "Random error.",
		},
		PropsCreator: create,
		Create: func(props *scaffui.SingleChildProps[ThrowerProps]) {
			props.OnLayout = func(node *scaffui.SingleChildNode[ThrowerProps]) (scath.Vec, error) {
				return scath.Vec{}, errors.New(node.Props().Msg)
			}
		},
	})
}
