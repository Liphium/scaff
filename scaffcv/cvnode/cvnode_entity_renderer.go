package cvnode

import (
	"maps"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
)

type EntityStore[K comparable, E any] struct {
	deletions []K
	updates   []K
	entities  map[K]E
}

func NewEntityStore[K comparable, E any]() *EntityStore[K, E] {
	return &EntityStore[K, E]{
		entities: map[K]E{},
	}
}

func (e *EntityStore[K, E]) Entity(key K, entity E) {
	if e.entities == nil {
		e.entities = make(map[K]E)
	}
	e.updates = append(e.updates, key)
	e.entities[key] = entity
}

func (e *EntityStore[K, E]) RemoveEntity(key K) {
	e.deletions = append(e.deletions, key)
	delete(e.entities, key)
}

func (e *EntityStore[K, E]) clearDiff() {
	e.deletions = []K{}
	e.updates = []K{}
}

type EntityRendererProps[K comparable, E any] struct {
	*EntityStore[K, E]
	Visible  bool
	Renderer func(key K, entity E) scaffcv.NodeBuilder
	scaffcv.AcceptNoChild
}

func (e EntityRendererProps[K, E]) render(key K, context *scaffcv.BuildContext) scaffcv.Node {
	if e.Renderer == nil {
		return nil
	}

	if en, ok := e.EntityStore.entities[key]; ok {
		return e.Renderer(key, en)(context)
	}
	return nil
}

func EntityRenderer[K comparable, E any](create func(t *scaff.Tracker, props *EntityRendererProps[K, E])) scaffcv.NodeBuilder {
	rendered := map[K]scaffcv.Node{}

	// Function to help render changed nodes and mount them properly
	renderChanged := func(parent scaffcv.Node, props EntityRendererProps[K, E], context *scaffcv.BuildContext) {

		// Unload all deleted nodes
		for _, key := range props.EntityStore.deletions {
			if node, ok := rendered[key]; ok {
				node.Unload()
				delete(rendered, key)
			}
		}

		// Mount all new nodes
		for _, key := range props.EntityStore.updates {
			if node, ok := rendered[key]; ok {
				node.Unload()
				delete(rendered, key)
			}

			node := props.render(key, context)
			if node != nil {
				node.Load(parent)
				rendered[key] = node
			}
		}

		props.EntityStore.clearDiff()
	}

	// Function to help calculate the minimum and maximum position (+ size) of all nodes
	calculateBounds := func() (min scath.Vec, max scath.Vec) {
		for _, node := range rendered {
			pos := node.Position()
			size := node.Size()

			if pos.X < min.X {
				min.X = pos.X
			}
			if pos.Y < min.Y {
				min.Y = pos.Y
			}
			if pos.X+size.X > max.X {
				max.X = pos.X + size.X
			}
			if pos.Y+size.Y > max.Y {
				max.Y = pos.Y + size.Y
			}
		}
		return min, max
	}

	return scaffcv.Standard(scaffcv.StandardCreate[EntityRendererProps[K, E]]{
		ID: "entity-renderer",
		DefaultProps: EntityRendererProps[K, E]{
			Visible: true,
			EntityStore: &EntityStore[K, E]{
				entities: map[K]E{},
			},
		},
		PropsCreator: create,
		Create: func(methods *scaffcv.StandardMethods[EntityRendererProps[K, E]]) {
			methods.Position = func(node *scaffcv.StandardNode[EntityRendererProps[K, E]]) scath.Vec {
				min, _ := calculateBounds()
				return min
			}
			methods.Size = func(node *scaffcv.StandardNode[EntityRendererProps[K, E]]) scath.Vec {
				min, max := calculateBounds()
				return max.Sub(min)
			}

			methods.OnPropsChanged = func(node *scaffcv.StandardNode[EntityRendererProps[K, E]]) {
				renderChanged(node, node.Props(), node.Context())
			}

			methods.OnUpdate = func(node *scaffcv.StandardNode[EntityRendererProps[K, E]], c *scaff.Context) error {
				renderChanged(node, node.Props(), node.Context())

				// Forward update to all nodes
				for _, node := range rendered {
					node.Update(c)
				}
				return nil
			}

			methods.OnDraw = func(node *scaffcv.StandardNode[EntityRendererProps[K, E]], c *scaff.Context, painter paint.Painter) {
				if node.Props().Visible {
					scaffcv.DrawCulled(maps.Values(rendered), c, painter, node.Context())
				}
			}

			methods.OnHandleEvent = func(node *scaffcv.StandardNode[EntityRendererProps[K, E]], c *scaff.Context, event scaff.Event) error {
				for _, node := range rendered {
					if err := node.HandleEvent(c, event); err != nil {
						return err
					}
				}
				return nil
			}
		},
	})
}
