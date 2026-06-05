package cvnode

import (
	"slices"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/hajimehoshi/ebiten/v2"
)

type KeyboardProps struct {
	// All of the keys to listen for
	Keys []ebiten.Key

	// Called when a key in the list is pressed (return if anything has been done with the event, will not call when the press has already been handled)
	OnPress func(key ebiten.Key) bool

	// Called when a key in the list is released (only called when the keypress was actually handled)
	OnRelease func(key ebiten.Key)

	// Called for every update where the key is pressed
	OnUpdate func(pressed []ebiten.Key) error

	*scaffcv.AcceptChild
}

func Keyboard(create func(t *scaff.Tracker, props *KeyboardProps)) scaffcv.NodeBuilder {
	return scaffcv.Standard(scaffcv.StandardCreate[KeyboardProps]{
		ID: "keyboard",
		DefaultProps: KeyboardProps{
			AcceptChild: &scaffcv.AcceptChild{},
		},
		PropsCreator: create,
		Create: func(methods *scaffcv.StandardMethods[KeyboardProps]) {
			pressed := []ebiten.Key{}

			methods.OnHandleEvent = func(node *scaffcv.StandardNode[KeyboardProps], c *scaff.Context, event scaff.Event) error {

				// Forward unhandled key presses
				if keyPress, ok := event.(scaff.KeyPressEvent); ok && !c.IsHandled(keyPress.EventID()) && slices.Contains(node.Props().Keys, keyPress.Key) {
					handle := node.Props().OnPress == nil || node.Props().OnPress(keyPress.Key)

					if handle {
						pressed = append(pressed, keyPress.Key)
						c.Handled(keyPress.EventID())
					}
				}

				// Handle key releases when the press was handled + call handler
				if keyRelease, ok := event.(scaff.KeyReleaseEvent); ok && !c.IsHandled(keyRelease.EventID()) && slices.Contains(node.Props().Keys, keyRelease.Key) {
					lenBefore := len(pressed)
					pressed = slices.DeleteFunc(pressed, func(k ebiten.Key) bool {
						return k == keyRelease.Key
					})

					// If the list doesn't have the same length anymore, we handle the release event
					if lenBefore != len(pressed) {
						if node.Props().OnRelease != nil {
							node.Props().OnRelease(keyRelease.Key)
						}
						c.Handled(keyRelease.EventID())
					}
				}

				return nil
			}

			// Call the update function with all pressed keys in case there are any
			methods.OnUpdate = func(node *scaffcv.StandardNode[KeyboardProps], c *scaff.Context) error {
				if len(pressed) > 0 && node.Props().OnUpdate != nil {
					return node.Props().OnUpdate(pressed)
				}
				return nil
			}
		},
	})
}
