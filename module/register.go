package module

import (
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-rendering/rendering"
	"github.com/adm87/finch-rendering/sprites"
	"github.com/adm87/finch-rendering/vector"
)

func RegisterModule() error {
	if err := rendering.RegisterRenderers(map[ecs.ComponentType]rendering.Renderer{
		sprites.SpriteRenderComponentType: sprites.SpriteRenderer,
		vector.BoxRenderComponentType:     vector.BoxRenderer,
	}); err != nil {
		return err
	}
	return nil
}
