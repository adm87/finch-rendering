package rendering

import (
	"slices"

	"github.com/adm87/finch-core/components/camera"
	"github.com/adm87/finch-core/ecs"
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

	queue, err := rs.GetRenderingQueue(world)
	if err != nil {
		return err
	}

	for _, pair := range queue {
		pair.Second(buffer, view)
	}
	return nil
}

func (rs *RenderingSystem) GetRenderingQueue(world *ecs.World) ([]types.Pair[int, RenderingTask], error) {
	var renderingQueue []types.Pair[int, RenderingTask]

	for ct := range rendererRegistry {
		entities := world.FilterEntitiesByComponents(ct)
		for entity := range entities {
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
