package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type slicer interface {
	slice() []rectF64
}

type collidable interface {
	slicer
	collEnabled() bool
}

type drawable interface {
	slicer
	drawEnabled() bool
	Color() color.Color
}

func draw(dst *ebiten.Image, src drawable) {
	if !src.drawEnabled() {
		return
	}

	for _, rect := range src.slice() {
		rect.draw(dst, src.Color())
	}
}

func collides(a, b collidable, tolerance float64) bool {
	if !a.collEnabled() || !b.collEnabled() {
		return false
	}

	for _, rectA := range a.slice() {
		for _, rectB := range b.slice() {
			if !intersects(rectA, rectB, tolerance) {
				continue
			}

			return true
		}
	}

	return false
}
