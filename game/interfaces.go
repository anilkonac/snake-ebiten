package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type collidable interface {
	isActive() bool
	rectSlice() []rectF64
}

type drawable interface {
	collidable
	Color() color.Color
}

func draw(src drawable, dst *ebiten.Image) {
	if !src.isActive() {
		return
	}

	for _, rect := range src.rectSlice() {
		ebitenutil.DrawRect(dst, rect.x, rect.y, rect.width, rect.height, src.Color())
		if debugRects {
			ebitenutil.DebugPrintAt(dst, fmt.Sprintf("%3.3f, %3.3f", rect.x, rect.y), int(rect.x)-90, int(rect.y)-15)
			bottomX := rect.x + rect.width
			bottomY := rect.y + rect.height
			ebitenutil.DebugPrintAt(dst, fmt.Sprintf("%3.3f, %3.3f", bottomX, bottomY), int(bottomX), int(bottomY))
		}
	}
}

func collides(a, b collidable) bool {
	if !a.isActive() || !b.isActive() {
		return false
	}

	for _, rectA := range a.rectSlice() {
		for _, rectB := range b.rectSlice() {
			if !intersects(rectA, rectB) {
				continue
			}

			return true
		}
	}

	return false
}
