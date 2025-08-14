package renderers

import "github.com/hajimehoshi/ebiten/v2"

// Renderer is an interface that defines a mechanism for rendering to an image.
type Renderer interface {
	Dispose()

	Render(buffer *ebiten.Image, view, transform ebiten.GeoM) error
}
