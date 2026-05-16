package scaffui

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
)

type ViewportProps struct {
	child optional.O[NodeBuilder]
}

func Viewport(create func(*scaff.Tracker)) scaff.NodeBuilder {
	return scaff.UseNode("viewport", func(props *scaff.SingleChildProps[ViewportProps]) {

	})
}
