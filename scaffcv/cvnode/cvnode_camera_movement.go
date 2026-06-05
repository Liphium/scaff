package cvnode

import (
	"maps"
	"slices"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type CameraMovementProps struct {
	Speed float64
	*scaffcv.AcceptChild
}

func CameraMovement(cameraPosition *scaff.Signal[scath.Vec], create func(t *scaff.Tracker, props *CameraMovementProps)) scaffcv.NodeBuilder {
	return scaffcv.Standard(scaffcv.StandardCreate[CameraMovementProps]{
		ID: "camera-movement",
		DefaultProps: CameraMovementProps{
			Speed:       10,
			AcceptChild: scaffcv.EmptyChild(),
		},
		PropsCreator: create,
		Create: func(props *scaffcv.StandardMethods[CameraMovementProps]) {
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

			props.OnPropsChanged = func(node *scaffcv.StandardNode[CameraMovementProps]) {
				node.Props().AcceptChild.Child(Keyboard(func(t *scaff.Tracker, props *KeyboardProps) {
					props.Keys = slices.Collect(maps.Keys(directionForKey))

					props.OnUpdate = func(pressed []ebiten.Key) error {
						active := []scath.Vec{}

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
				}))
			}
		},
	})
}
