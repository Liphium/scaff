package scaffui

import (
	"github.com/Liphium/scaff"
)

// Wrap templates for use in scaffui

type AcceptNoChild = scaff.AcceptNoChildTemplate[NodeBuilder]
type AcceptChild = scaff.AcceptChildTemplate[NodeBuilder]
type AcceptChildren = scaff.AcceptChildrenTemplate[NodeBuilder]
