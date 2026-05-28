package scaff

type RootProps struct {
	children []NodeBuilder
}

func (sp *RootProps) Child(builder NodeBuilder) {
	sp.children = append(sp.children, builder)
}

func (sp RootProps) GetBuilders() []NodeBuilder {
	return sp.children
}

	func root(create func(t *Tracker, props *RootProps)) NodeBuilder {
		return MultiNode(MultiNodeCreate[RootProps]{
			ID:           "root",
			PropsCreator: create,
		})
	}
