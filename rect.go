package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type screenRect struct {
	x, y   float64
	width  float64
	height float64
}

func newScreenRect(x, y, width, height float64) *screenRect {
	return &screenRect{x: x, y: y, width: width, height: height}
}

func (r screenRect) draw(dst *ebiten.Image, color color.Color) {
	ebitenutil.DrawRect(dst, r.x, r.y, r.width, r.height, color)
}
