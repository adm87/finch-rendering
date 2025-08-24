package vector

import (
	"image/color"

	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/geometry"
)

var BoxRenderComponentType = ecs.NewComponentType[*BoxRenderComponent]()

type BoxRenderComponent struct {
	Min         geometry.Point64
	Max         geometry.Point64
	IsLocal     bool
	DrawBorder  bool
	DrawFill    bool
	BorderColor color.RGBA
	FillColor   color.RGBA
	BorderWidth float32
	ZOrder      int
}

func NewBoxRenderComponent(zOrder int) *BoxRenderComponent {
	return &BoxRenderComponent{
		DrawFill:    true,
		BorderColor: color.RGBA{0, 0, 0, 255},
		FillColor:   color.RGBA{255, 255, 255, 255},
		ZOrder:      zOrder,
		BorderWidth: 1,
	}
}

func (brc *BoxRenderComponent) Type() ecs.ComponentType {
	return BoxRenderComponentType
}
