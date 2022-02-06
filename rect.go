package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type screenRect struct {
	x, y          float64
	width, height float64
}

func (r screenRect) intersects(other *screenRect) bool {

	// intersectsX := (other.x > r.x && other.x < r.x+r.width) || (r.x > other.x && r.x < other.x+other.width)
	// intersectsY := (other.y > r.y && other.y < r.y+r.height) || (r.y > other.y && r.y < other.y+other.height)

	// return intersectsX && intersectsY

	return false
}

func (r screenRect) draw(dst *ebiten.Image, clr color.Color) {
	drawRect(dst, r.x, r.y, r.width, r.height, clr)
}
