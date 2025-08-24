package sprites

import (
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/geometry"
	"github.com/hajimehoshi/ebiten/v2"
)

var SpriteRenderComponentType = ecs.NewComponentType[*SpriteRenderComponent]()

type SpriteRenderComponent struct {
	ImageID    string
	ZOrder     int
	Blend      ebiten.Blend
	ColorScale ebiten.ColorScale
	Filter     ebiten.Filter
	Anchor     geometry.Point64
}

func NewSpriteRenderComponent(imageID string, zOrder int) *SpriteRenderComponent {
	return &SpriteRenderComponent{
		ImageID: imageID,
		ZOrder:  zOrder,
	}
}

func (src *SpriteRenderComponent) Type() ecs.ComponentType {
	return SpriteRenderComponentType
}
