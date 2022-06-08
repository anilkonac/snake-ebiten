package core

import (
	"image/color"

	t "github.com/anilkonac/snake-ebiten/game/tool"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Rectangle compatible with float32 type parameters of the ebiten.DrawTriangleShader function.
type RectF32 struct {
	Pos       t.Vec32
	Size      t.Vec32
	PosInUnit t.Vec32
}

func NewRect(pos, size t.Vec32) *RectF32 {
	return &RectF32{pos, size, t.Vec32{X: 0, Y: 0}}
}

func (r RectF32) DrawOuterRect(dst *ebiten.Image, clr color.Color) {
	pos64 := r.Pos.To64()
	size64 := r.Size.To64()
	ebitenutil.DrawRect(dst, pos64.X, pos64.Y, size64.X, size64.Y, color.RGBA{255, 255, 255, 96})
}
