package sprites

import (
	"github.com/adm87/finch-core/components/transform"
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/geometry"
	"github.com/adm87/finch-rendering/rendering"
	"github.com/adm87/finch-resources/images"
	"github.com/hajimehoshi/ebiten/v2"
)

var op = &ebiten.DrawImageOptions{}

func SpriteRenderer(world *ecs.World, entity ecs.Entity) (rendering.RenderingTask, int, error) {
	spriteComp, _, _ := ecs.GetComponent[*SpriteRenderComponent](world, entity, SpriteRenderComponentType)
	transformComp, _, _ := ecs.GetComponent[*transform.TransformComponent](world, entity, transform.TransformComponentType)
	img, err := images.Cache().Get(spriteComp.ImageID)
	if err != nil {
		return nil, 0, err
	}
	origin := geometry.Point64{
		X: spriteComp.Anchor.X * float64(img.Bounds().Dx()),
		Y: spriteComp.Anchor.Y * float64(img.Bounds().Dy()),
	}
	return func(surface *ebiten.Image, view ebiten.GeoM) {
		op.GeoM.Reset()
		op.GeoM.Translate(-origin.X, -origin.Y)
		op.GeoM.Concat(transformComp.WorldMatrix())
		op.GeoM.Concat(view)

		op.Blend = spriteComp.Blend
		op.ColorScale = spriteComp.ColorScale
		op.Filter = spriteComp.Filter

		surface.DrawImage(img, op)
	}, spriteComp.ZOrder, nil
}
