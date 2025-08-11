package rendering

import (
	"slices"

	"github.com/adm87/finch-core/components/transform"
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/hash"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	RenderSystemType   = ecs.SystemType(hash.GetHashFromType[RenderSystem]())
	RenderSystemFilter = []ecs.ComponentType{
		transform.TransformComponentType,
		RenderComponentType,
	}
)

type RenderQueueItem struct {
	Renderer  Renderer
	Order     int
	Transform ebiten.GeoM
}

type RenderSystem struct {
}

func NewRenderSystem() *RenderSystem {
	return &RenderSystem{}
}

func (s *RenderSystem) Type() ecs.SystemType {
	return RenderSystemType
}

func (s *RenderSystem) Filter() []ecs.ComponentType {
	return RenderSystemFilter
}

func (s *RenderSystem) Render(entities hash.HashSet[ecs.Entity], buffer *ebiten.Image, view ebiten.GeoM) error {
	queue, err := internal_get_render_queue(entities)
	if err != nil {
		return err
	}

	for _, item := range queue {
		if err := item.Renderer.Render(buffer, view, item.Transform); err != nil {
			return err
		}
	}

	return nil
}

func internal_get_render_queue(entities hash.HashSet[ecs.Entity]) ([]*RenderQueueItem, error) {
	var items []*RenderQueueItem

	for entity := range entities {
		rc, _, rErr := ecs.GetComponent[*RenderComponent](entity, RenderComponentType)
		if rErr != nil {
			return nil, rErr
		}

		if !rc.IsVisible || rc.Renderer == nil {
			continue
		}

		tc, _, tErr := ecs.GetComponent[*transform.TransformComponent](entity, transform.TransformComponentType)
		if tErr != nil {
			return nil, tErr
		}

		items = append(items, &RenderQueueItem{
			Renderer:  rc.Renderer,
			Order:     rc.ZIndex,
			Transform: tc.WorldMatrix(),
		})
	}

	slices.SortFunc(items, func(a, b *RenderQueueItem) int {
		return a.Order - b.Order
	})

	return items, nil
}
