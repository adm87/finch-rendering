package vector

import (
	"image/color"
	"math"

	"github.com/adm87/finch-core/geometry"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var lineImg *ebiten.Image

func init() {
	lineImg = ebiten.NewImage(1, 1)
	lineImg.Fill(color.White)
}

type BoxRenderer struct {
	DrawBorder bool
	DrawFill   bool

	min geometry.Point64
	max geometry.Point64

	border color.RGBA
	fill   color.RGBA
}

func NewBoxRenderer() *BoxRenderer {
	return &BoxRenderer{
		border: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		fill:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
	}
}

func (b *BoxRenderer) SetArea(min geometry.Point64, max geometry.Point64) {
	b.min = min
	b.max = max
}

func (b *BoxRenderer) SetBorder(color color.RGBA) {
	b.border = color
}

func (b *BoxRenderer) SetFill(color color.RGBA) {
	b.fill = color
}

func (r *BoxRenderer) Dispose() {

}

func (r *BoxRenderer) Render(buffer *ebiten.Image, view, transform ebiten.GeoM) error {
	path := vector.Path{}

	minX := math.Min(r.min.X, r.max.X)
	minY := math.Min(r.min.Y, r.max.Y)
	maxX := math.Max(r.min.X, r.max.X)
	maxY := math.Max(r.min.Y, r.max.Y)

	minX, minY = view.Apply(float64(minX), float64(minY))
	maxX, maxY = view.Apply(float64(maxX), float64(maxY))

	path.MoveTo(float32(minX), float32(minY))
	path.LineTo(float32(maxX), float32(minY))
	path.LineTo(float32(maxX), float32(maxY))
	path.LineTo(float32(minX), float32(maxY))
	path.Close()

	if r.DrawFill {
		r.RenderFill(path, buffer, &ebiten.DrawTrianglesOptions{})
	}

	if r.DrawBorder {
		r.RenderBorder(path, buffer, &ebiten.DrawTrianglesOptions{})
	}

	return nil
}

func (r *BoxRenderer) RenderBorder(path vector.Path, buffer *ebiten.Image, op *ebiten.DrawTrianglesOptions) {
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: 1,
	})

	if len(vs) == 0 || len(is) == 0 {
		return
	}

	for i := range vs {
		vs[i].ColorR = float32(r.border.R) / 255.0
		vs[i].ColorG = float32(r.border.G) / 255.0
		vs[i].ColorB = float32(r.border.B) / 255.0
		vs[i].ColorA = float32(r.border.A) / 255.0
	}

	buffer.DrawTriangles(vs, is, lineImg, op)
}

func (r *BoxRenderer) RenderFill(path vector.Path, buffer *ebiten.Image, op *ebiten.DrawTrianglesOptions) {
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)

	if len(vs) == 0 || len(is) == 0 {
		return
	}

	for i := range vs {
		vs[i].ColorR = float32(r.fill.R) / 255.0
		vs[i].ColorG = float32(r.fill.G) / 255.0
		vs[i].ColorB = float32(r.fill.B) / 255.0
		vs[i].ColorA = float32(r.fill.A) / 255.0
	}

	buffer.DrawTriangles(vs, is, lineImg, op)
}
