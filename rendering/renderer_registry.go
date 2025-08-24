package rendering

import "github.com/adm87/finch-core/ecs"

var rendererRegistry = make(map[ecs.ComponentType]Renderer)

func Register(componentType ecs.ComponentType, renderer Renderer) {
	rendererRegistry[componentType] = renderer
}
