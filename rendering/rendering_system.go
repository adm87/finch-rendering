package rendering

import (
	"slices"

	"github.com/adm87/finch-core/components/camera"
	"github.com/adm87/finch-core/components/transform"
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/geometry"
	"github.com/adm87/finch-core/types"
	"github.com/hajimehoshi/ebiten/v2"
)

var RenderingSystemType = ecs.NewSystemType[*RenderingSystem]()

type RenderingSystem struct {
	enabled bool
}

type (
	// RenderingTask is a function that performs the rendering for a specific component.
	RenderingTask func(surface *ebiten.Image, view ebiten.GeoM)
	// Renderer is a function that creates the RenderingTask for an entity.
	//
	// It returns the RenderingTask, the Z-order for the entity, and any error that occurred.
	Renderer func(world *ecs.World, entity ecs.Entity) (RenderingTask, int, error)
)

func NewRenderingSystem() *RenderingSystem {
	return &RenderingSystem{
		enabled: true,
	}
}

func (rs *RenderingSystem) IsEnabled() bool {
	return rs.enabled
}

func (rs *RenderingSystem) Enable() {
	rs.enabled = true
}

func (rs *RenderingSystem) Disable() {
	rs.enabled = false
}

func (rs *RenderingSystem) Type() ecs.SystemType {
	return RenderingSystemType
}

func (rs *RenderingSystem) Render(world *ecs.World, buffer *ebiten.Image) error {
	cameraComp, err := camera.FindCameraComponent(world)
	if err != nil {
		return err
	}

	view := cameraComp.WorldMatrix()
	view.Invert()

	viewport := cameraComp.Viewport()

	queue, err := rs.GetRenderingQueue(world, viewport)
	if err != nil {
		return err
	}

	for _, pair := range queue {
		pair.Second(buffer, view)
	}
	return nil
}

func (rs *RenderingSystem) GetRenderingQueue(world *ecs.World, viewport geometry.Rectangle64) ([]types.Pair[int, RenderingTask], error) {
	var renderingQueue []types.Pair[int, RenderingTask]

	for ct := range rendererRegistry {
		entities := world.FilterEntitiesByComponents(ct)
		for entity := range entities {
			if !is_visible(world, entity, viewport) {
				continue
			}
			renderingTask, zOrder, err := rendererRegistry[ct](world, entity)
			if err != nil {
				return nil, err
			}
			renderingQueue = append(renderingQueue, types.Pair[int, RenderingTask]{First: zOrder, Second: renderingTask})
		}
	}

	slices.SortFunc(renderingQueue, func(a, b types.Pair[int, RenderingTask]) int {
		return a.First - b.First
	})

	return renderingQueue, nil
}

func is_visible(world *ecs.World, entity ecs.Entity, viewport geometry.Rectangle64) bool {
	visibilityComp, found, _ := ecs.GetComponent[*VisibilityComponent](world, entity, VisibilityComponentType)
	if !found {
		return true // No visibility component present so the entity is always visible.
	}
	if !visibilityComp.IsVisible {
		return false // IsVisible takes priority over anything.
	}
	if !visibilityComp.VisibleArea.IsValid() {
		return true // Visible area is not defined, but the entity is marked as visible so we assume it's visible
	}
	if transform, found, _ := ecs.GetComponent[*transform.TransformComponent](world, entity, transform.TransformComponentType); found {
		position := transform.Position()
		visibleArea := visibilityComp.VisibleArea.Value()
		return visibleArea.Translate(position.X, position.Y).Intersects(viewport)
	}
	return false // Warning: If we can't find a transform to translate to world space, we assume it's not visible
}
