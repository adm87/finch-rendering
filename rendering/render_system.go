package rendering

import (
	"slices"

	"github.com/adm87/finch-core/components/camera"
	"github.com/adm87/finch-core/components/transform"
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-rendering/renderers"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	RenderSystemType   = ecs.NewSystemType[*RenderSystem]()
	RenderSystemFilter = []ecs.ComponentType{
		transform.TransformComponentType,
		RenderComponentType,
	}
)

type RenderQueueItem struct {
	Renderer  renderers.Renderer
	Order     int
	Transform ebiten.GeoM
}

type RenderSystem struct {
	enabled bool
}

func NewRenderSystem() *RenderSystem {
	return &RenderSystem{
		enabled: true,
	}
}

func (s *RenderSystem) Enable() {
	s.enabled = true
}

func (s *RenderSystem) Disable() {
	s.enabled = false
}

func (s *RenderSystem) IsEnabled() bool {
	return s.enabled
}

func (s *RenderSystem) Type() ecs.SystemType {
	return RenderSystemType
}

func (s *RenderSystem) Render(world *ecs.World, buffer *ebiten.Image) error {
	queue, err := internal_get_render_queue(world)
	if err != nil {
		return err
	}

	cameraComponent, err := camera.FindCameraComponent(world)
	if err != nil {
		return err
	}

	view := ebiten.GeoM{}
	if cameraComponent != nil {
		view = cameraComponent.WorldMatrix()
		view.Invert()
	}

	for _, item := range queue {
		if err := item.Renderer.Render(buffer, view, item.Transform); err != nil {
			return err
		}
	}

	return nil
}

func internal_get_render_queue(world *ecs.World) ([]*RenderQueueItem, error) {
	entities := world.FilterEntitiesByComponents(RenderSystemFilter...)

	var items []*RenderQueueItem

	for entity := range entities {
		rc, _, rErr := ecs.GetComponent[*RenderComponent](world, entity, RenderComponentType)
		if rErr != nil {
			return nil, rErr
		}

		if !rc.IsVisible || rc.Renderer == nil {
			continue
		}

		tc, _, tErr := ecs.GetComponent[*transform.TransformComponent](world, entity, transform.TransformComponentType)
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
