package vector

import (
	"image/color"
	"math"

	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-core/geometry"
	"github.com/adm87/finch-rendering/rendering"
	"github.com/hajimehoshi/ebiten/v2"
	vec "github.com/hajimehoshi/ebiten/v2/vector"
)

var boxImg = ebiten.NewImage(3, 3)
var boxOp = &ebiten.DrawTrianglesOptions{}

func BoxRenderer(world *ecs.World, entity ecs.Entity) (rendering.RenderingTask, int, error) {
	boxComp, _, _ := ecs.GetComponent[*BoxRenderComponent](world, entity, BoxRenderComponentType)
	left, right, top, bottom := get_edges(boxComp.Min, boxComp.Max)
	return func(surface *ebiten.Image, view ebiten.GeoM) {
		if !boxComp.DrawBorder && !boxComp.DrawFill {
			return
		}

		path := vec.Path{}

		left, top = view.Apply(left, top)
		right, bottom = view.Apply(right, bottom)

		path.MoveTo(float32(left), float32(top))
		path.LineTo(float32(right), float32(top))
		path.LineTo(float32(right), float32(bottom))
		path.LineTo(float32(left), float32(bottom))
		path.Close()

		if boxComp.DrawBorder {
			r, g, b, a := float32(boxComp.BorderColor.R)/255, float32(boxComp.BorderColor.G)/255, float32(boxComp.BorderColor.B)/255, float32(boxComp.BorderColor.A)/255
			draw_border(&path, surface, []float32{r, g, b, a}, boxComp.BorderWidth, boxOp)
		}

		if boxComp.DrawFill {
			r, g, b, a := float32(boxComp.FillColor.R)/255, float32(boxComp.FillColor.G)/255, float32(boxComp.FillColor.B)/255, float32(boxComp.FillColor.A)/255
			draw_fill(&path, surface, []float32{r, g, b, a}, boxOp)
		}

	}, boxComp.ZOrder, nil
}

func draw_border(path *vec.Path, surface *ebiten.Image, color []float32, width float32, op *ebiten.DrawTrianglesOptions) {
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, &vec.StrokeOptions{
		Width: width,
	})
	if len(vs) == 0 || len(is) == 0 {
		return
	}
	for i := range vs {
		vs[i].ColorR = color[0]
		vs[i].ColorG = color[1]
		vs[i].ColorB = color[2]
		vs[i].ColorA = color[3]
	}
	surface.DrawTriangles(vs, is, boxImg, op)
}

func draw_fill(path *vec.Path, surface *ebiten.Image, color []float32, op *ebiten.DrawTrianglesOptions) {
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	if len(vs) == 0 || len(is) == 0 {
		return
	}
	for i := range vs {
		vs[i].ColorR = color[0]
		vs[i].ColorG = color[1]
		vs[i].ColorB = color[2]
		vs[i].ColorA = color[3]
	}
	surface.DrawTriangles(vs, is, boxImg, op)
}

func get_edges(min, max geometry.Point64) (left, right, top, bottom float64) {
	left = math.Min(min.X, max.X)
	right = math.Max(min.X, max.X)
	top = math.Min(min.Y, max.Y)
	bottom = math.Max(min.Y, max.Y)
	return
}

func init() {
	boxImg.Fill(color.White)
	rendering.Register(BoxRenderComponentType, BoxRenderer)
}
