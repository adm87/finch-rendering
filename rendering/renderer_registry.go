package rendering

import (
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/errors"
)

var rendererRegistry = make(map[ecs.ComponentType]Renderer)

func RegisterRenderers(renderers map[ecs.ComponentType]Renderer) error {
	for componentType, renderer := range renderers {
		if _, exists := rendererRegistry[componentType]; exists {
			return errors.NewDuplicateError("renderer already registered for component type: " + componentType.String())
		}
		rendererRegistry[componentType] = renderer
	}
	return nil
}
