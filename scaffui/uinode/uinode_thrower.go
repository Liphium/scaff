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
	return scaffui.Standard(scaffui.StandardCreate[ThrowerProps]{
		ID: "thrower",
		DefaultProps: ThrowerProps{
			AcceptNoChild: &scaffui.AcceptNoChild{},
			Msg:           "Random error.",
		},
		PropsCreator: create,
		Create: func(props *scaffui.StandardMethods[ThrowerProps]) {
			props.OnLayout = func(node *scaffui.StandardNode[ThrowerProps]) (scath.Vec, error) {
				return scath.Vec{}, errors.New(node.Props().Msg)
			}
		},
	})
}
