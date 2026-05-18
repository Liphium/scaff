package main

import "github.com/Liphium/scaff"

type StackProps struct {
	children []scaff.NodeBuilder
}

func (sp *StackProps) Child(builder scaff.NodeBuilder) {
	sp.children = append(sp.children, builder)
}

func (sp StackProps) GetBuilders() []scaff.NodeBuilder {
	return sp.children
}

func Stack(create func(t *scaff.Tracker, props *StackProps)) scaff.NodeBuilder {
	return scaff.CreateMultiNode("stack", create, nil)
}
