package rendering

import (
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/geometry"
	"github.com/adm87/finch-core/types"
)

var VisibilityComponentType = ecs.NewComponentType[*VisibilityComponent]()

type VisibilityComponent struct {
	IsVisible   bool
	VisibleArea types.Optional[geometry.Rectangle64]
}

// NewVisibilityComponent creates a new VisibilityComponent to control the visibility of an entity within the rendering system.
//
// Entities with a render component but not a visibility component are always considered visible. The visible area is in local coordinates.
func NewVisibilityComponent() *VisibilityComponent {
	return &VisibilityComponent{
		IsVisible:   true,
		VisibleArea: types.NewEmptyOptional[geometry.Rectangle64](),
	}
}

func (c *VisibilityComponent) Type() ecs.ComponentType {
	return VisibilityComponentType
}
