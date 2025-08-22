package sprites

import (
	"github.com/adm87/finch-core/geometry"
	"github.com/adm87/finch-core/types"
	"github.com/adm87/finch-resources/storage"
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteRenderer struct {
	image  storage.ImageHandle
	anchor geometry.Point64
	origin types.Optional[geometry.Point64]
}

func NewSpriteRenderer(texName string, anchor geometry.Point64) *SpriteRenderer {
	return &SpriteRenderer{
		image:  storage.ImageHandle(texName),
		anchor: anchor,
		origin: types.NewEmptyOptional[geometry.Point64](),
	}
}

func (r *SpriteRenderer) Size() (size geometry.Point64) {
	img, err := r.image.Get()
	if err != nil {
		panic(err)
	}
	r.origin.Invalidate()
	return geometry.Point64{
		X: float64(img.Bounds().Dx()),
		Y: float64(img.Bounds().Dy()),
	}
}

func (r *SpriteRenderer) Dispose() {

}

func (r *SpriteRenderer) Render(buffer *ebiten.Image, view, transform ebiten.GeoM) error {
	img, err := r.image.Get()
	if err != nil {
		return err
	}

	if !r.origin.IsValid() {
		r.origin.SetValue(geometry.Point64{
			X: float64(img.Bounds().Dx()) * r.anchor.X,
			Y: float64(img.Bounds().Dy()) * r.anchor.Y,
		})
	}
	origin := r.origin.Value()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-origin.X, -origin.Y)
	op.GeoM.Concat(transform)
	op.GeoM.Concat(view)
	buffer.DrawImage(img, op)
	return nil
}
