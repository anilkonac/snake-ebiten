package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type screenRect struct {
	x, y          float64
	width, height float64
}

func (a screenRect) intersects(b *screenRect) bool {

	aBottomRightX := a.x + a.width
	aBottomRightY := a.y + a.height
	bBottomRightX := b.x + b.width
	bBottomRightY := b.y + b.height

	if (a.x < b.x) && (aBottomRightX < b.x) { // a is on the left side of b
		return false
	}

	if (a.x > bBottomRightX) && (aBottomRightX > bBottomRightX) { // a is on the right side of b
		return false
	}

	if (a.y < b.y) && (aBottomRightY < b.y) { // a is above b
		return false
	}

	if (a.y > bBottomRightY) && (aBottomRightY > bBottomRightY) { // a is under b
		return false
	}

	return true
}

func (r screenRect) draw(dst *ebiten.Image, clr color.Color) {
	drawRect(dst, r.x, r.y, r.width, r.height, clr)
}
