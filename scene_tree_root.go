package scaff

type RootProps struct {
	*AcceptChildren
}

func root(create func(t *Tracker, props *RootProps)) NodeBuilder {
	return Standard(StandardCreate[RootProps]{
		ID: "root",
		DefaultProps: RootProps{
			AcceptChildren: &AcceptChildren{},
		},
		PropsCreator: create,
	})
}
