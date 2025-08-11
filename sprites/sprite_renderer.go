package sprites

import (
	"github.com/adm87/finch-core/geometry"
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteRenderer struct {
	image  *ebiten.Image
	origin geometry.Point64
}

func NewSpriteRenderer(image *ebiten.Image, anchor geometry.Point64) *SpriteRenderer {
	origin := geometry.Point64{
		X: anchor.X * float64(image.Bounds().Dx()),
		Y: anchor.Y * float64(image.Bounds().Dy()),
	}
	return &SpriteRenderer{
		image:  image,
		origin: origin,
	}
}

func (r *SpriteRenderer) Dispose() {
	r.image = nil
}

func (r *SpriteRenderer) Render(buffer *ebiten.Image, view, transform ebiten.GeoM) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-r.origin.X, -r.origin.Y)
	op.GeoM.Concat(transform)
	op.GeoM.Concat(view)
	buffer.DrawImage(r.image, op)
	return nil
}
