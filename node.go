package scaff

import "github.com/hajimehoshi/ebiten/v2"

type Node interface {
	Tracking
	Identifiable

	// Should return your own parent
	Parent() Node

	// Should return your own children
	Children() []Node

	// Called when the Node is initialized in the tree
	Load(parent Node)

	// Called when the Node is removed from the tree
	Unload()

	// Called on every physics tick (like 60 times a second, depending on what ebitens tick rate is)
	Update(c *Context) *TracedError

	// Handle events from the system (you do not have to handle any, but should always push them along to children at least)
	HandleEvent(c *Context, event Event) *TracedError

	// Draw the thing onto the screen
	Draw(c *Context, image *ebiten.Image)
}
