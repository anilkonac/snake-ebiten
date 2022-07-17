package core

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Rectangle compatible with float32 type parameters of the ebiten.DrawTriangleShader function.
type RectF32 struct {
	Pos       Vec32
	Size      Vec32
	PosInUnit Vec32
}

func NewRect(pos, size Vec32) *RectF32 {
	return &RectF32{pos, size, Vec32{0, 0}}
}

func (r RectF32) DrawOuterRect(dst *ebiten.Image, clr color.Color) {
	pos64 := r.Pos.To64()
	size64 := r.Size.To64()
	ebitenutil.DrawRect(dst, pos64.X, pos64.Y, size64.X, size64.Y, color.RGBA{255, 255, 255, 96})
}

func MarkPoint(dst *ebiten.Image, p Vec64, length float64, clr color.Color) {
	ebitenutil.DrawLine(dst, p.X-length, p.Y, p.X+length, p.Y, clr)
	ebitenutil.DrawLine(dst, p.X, p.Y-length, p.X, p.Y+length, clr)
}
