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

func (r screenRect) draw(dst *ebiten.Image, color color.Color) {
	// Check if this rectangle is out of the screen
	// -----------------------------------
	if r.x < 0 { // left part is off the screen
		screenRect{screenWidth + r.x, r.y, -r.x, r.height}.draw(dst, color) // teleported left part
		screenRect{0, r.y, r.width + r.x, r.height}.draw(dst, color)        // part that is in the screen
		return
	} else if r.x+r.width > screenWidth { // right part is off the screen
		screenRect{0, r.y, r.x + r.width - screenWidth, r.height}.draw(dst, color) // teleported right part
		screenRect{r.x, r.y, screenWidth - r.x, r.height}.draw(dst, color)         // part that is in the screen
		return
	}

	if r.y < 0 { // upper part is off the screen
		screenRect{r.x, screenHeight + r.y, r.width, -r.y}.draw(dst, color) // teleported upper part
		screenRect{r.x, 0, r.width, r.height + r.y}.draw(dst, color)        // part that is in the screen
		return
	} else if r.y+r.height > screenHeight { // bottom part is off the screen
		screenRect{r.x, 0, r.width, r.y + r.height - screenHeight}.draw(dst, color) // teleported right part
		screenRect{r.x, r.y, r.width, screenHeight - r.y}.draw(dst, color)          // part that is in the screen
		return
	}

	// Draw the rectangle at last
	// --------------------------
	ebitenutil.DrawRect(dst, r.x, r.y, r.width, r.height, color)

}
