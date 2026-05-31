package cvnode

import (
	"slices"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type CameraMovementProps struct {
	Speed float64
	scaffcv.AcceptNoChild
}

func CameraMovement(cameraPosition *scaff.Signal[scath.Vec], create func(t *scaff.Tracker, props *CameraMovementProps)) scaffcv.NodeBuilder {
	return scaffcv.SingleNode(scaffcv.SingleNodeCreate[CameraMovementProps]{
		ID: "camera-movement",
		DefaultProps: CameraMovementProps{
			Speed: 10,
		},
		PropsCreator: create,
		Create: func(props *scaffcv.SingleChildProps[CameraMovementProps]) {
			directionForKey := map[ebiten.Key]scath.Vec{
				ebiten.KeyArrowUp: scath.Up,
				ebiten.KeyW:       scath.Up,

				ebiten.KeyArrowDown: scath.Down,
				ebiten.KeyS:         scath.Down,

				ebiten.KeyArrowRight: scath.Right,
				ebiten.KeyD:          scath.Right,

				ebiten.KeyArrowLeft: scath.Left,
				ebiten.KeyA:         scath.Left,
			}

			props.OnUpdate = func(node *scaffcv.SingleChildNode[CameraMovementProps], c *scaff.Context) error {
				active := []scath.Vec{}

				// TODO: Convert to scaff input API with events and stuff
				for key, vec := range directionForKey {
					if ebiten.IsKeyPressed(key) && !slices.ContainsFunc(active, func(v scath.Vec) bool {
						return v.Equals(vec)
					}) {
						active = append(active, vec)
					}
				}

				// Sum movement vectors + normalize
				sum := scath.Vec{}
				for _, v := range active {
					sum = sum.Add(v)
				}
				sum = sum.Unit()

				// Actually move the camera
				sum = sum.Scale(node.Props().Speed)
				cameraPosition.Set(cameraPosition.Value().Add(sum))
				return nil
			}
		},
	})
}
