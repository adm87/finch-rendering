package rendering

import (
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/errors"
)

var RenderComponentType = ecs.NewComponentType[*RenderComponent]()

type RenderComponent struct {
	Renderer  Renderer
	ZIndex    int
	IsVisible bool
}

func NewRenderComponent(renderer Renderer, zIndex int) *RenderComponent {
	if renderer == nil {
		panic(errors.NewInvalidArgumentError("renderer cannot be nil"))
	}
	return &RenderComponent{
		Renderer:  renderer,
		ZIndex:    zIndex,
		IsVisible: true,
	}
}

func (c *RenderComponent) Dispose() {
	c.Renderer = nil
}

func (c *RenderComponent) Type() ecs.ComponentType {
	return RenderComponentType
}
