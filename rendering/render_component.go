package rendering

import (
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/errors"
	"github.com/adm87/finch-core/hash"
)

var RenderComponentType = ecs.ComponentType(hash.GetHashFromType[RenderComponent]())

type RenderComponent struct {
	renderer  Renderer
	zIndex    int
	isVisible bool
}

func NewRenderComponent(renderer Renderer, zIndex int) *RenderComponent {
	if renderer == nil {
		panic(errors.NewInvalidArgumentError("renderer cannot be nil"))
	}
	return &RenderComponent{
		renderer:  renderer,
		zIndex:    zIndex,
		isVisible: true,
	}
}

func (c *RenderComponent) Type() ecs.ComponentType {
	return RenderComponentType
}

func (c *RenderComponent) Renderer() Renderer {
	return c.renderer
}

func (c *RenderComponent) SetRenderer(r Renderer) {
	if r == nil {
		panic(errors.NewInvalidArgumentError("renderer cannot be nil"))
	}
	c.renderer = r
}

func (c *RenderComponent) ZIndex() int {
	return c.zIndex
}

func (c *RenderComponent) SetZIndex(z int) {
	c.zIndex = z
}

func (c *RenderComponent) IsVisible() bool {
	return c.isVisible
}

func (c *RenderComponent) SetVisible(v bool) {
	c.isVisible = v
}
