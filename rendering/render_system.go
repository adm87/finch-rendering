package rendering

import (
	"slices"

	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/hash"
	"github.com/adm87/finch-core/transform"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	RenderSystemType   = ecs.SystemType(hash.GetHashFromType[RenderSystem]())
	RenderSystemFilter = []ecs.ComponentType{
		transform.TransformComponentType,
		RenderComponentType,
	}
)

type RenderSystem struct {
	cachedRenderComponents    map[ecs.EntityID]*RenderComponent
	cachedTransformComponents map[ecs.EntityID]*transform.TransformComponent
}

func NewRenderSystem() *RenderSystem {
	return &RenderSystem{
		cachedRenderComponents:    make(map[ecs.EntityID]*RenderComponent),
		cachedTransformComponents: make(map[ecs.EntityID]*transform.TransformComponent),
	}
}

func (s *RenderSystem) Type() ecs.SystemType {
	return RenderSystemType
}

func (s *RenderSystem) Filter() []ecs.ComponentType {
	return RenderSystemFilter
}

func (s *RenderSystem) Render(entities []*ecs.Entity, buffer *ebiten.Image, view ebiten.GeoM, interpolation float64) error {
	s.cachedRenderComponents = make(map[ecs.EntityID]*RenderComponent)
	s.cachedTransformComponents = make(map[ecs.EntityID]*transform.TransformComponent)

	renderOrder := s.internal_setup_rendering(entities)
	for _, entity := range renderOrder {
		rc := s.cachedRenderComponents[entity.ID()]
		tc := s.cachedTransformComponents[entity.ID()]

		if err := rc.Renderer().Render(buffer, view, tc.WorldMatrix()); err != nil {
			return err
		}
	}
	return nil
}

func (s *RenderSystem) internal_setup_rendering(entities []*ecs.Entity) []*ecs.Entity {
	var renderComponents []*ecs.Entity

	for _, entity := range entities {
		tc, rc, err := s.internal_get_operation_components(entity)
		if err != nil {
			continue
		}

		if !rc.IsVisible() {
			continue
		}

		s.cachedRenderComponents[entity.ID()] = rc
		s.cachedTransformComponents[entity.ID()] = tc

		renderComponents = append(renderComponents, entity)
	}

	slices.SortStableFunc(renderComponents, func(a, b *ecs.Entity) int {
		return s.cachedRenderComponents[a.ID()].ZIndex() - s.cachedRenderComponents[b.ID()].ZIndex()
	})

	return renderComponents
}

func (s *RenderSystem) internal_get_operation_components(entity *ecs.Entity) (*transform.TransformComponent, *RenderComponent, error) {
	tc, _, err := entity.GetComponent(transform.TransformComponentType)
	if err != nil {
		return nil, nil, err
	}
	rc, _, err := entity.GetComponent(RenderComponentType)
	if err != nil {
		return nil, nil, err
	}
	return tc.(*transform.TransformComponent), rc.(*RenderComponent), nil
}
