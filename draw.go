package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Draws given rectangle. If some part of the rectangle is off-screen, draw that part on the other edge.
func drawRect(dst *ebiten.Image, x, y, width, height float64, color color.Color) {
	// Check if the given rectangle is off-screen
	// ------------------------------------------
	if x < 0 { // left part is off-screen
		drawRect(dst, screenWidth+x, y, -x, height, color) // teleported left part
		drawRect(dst, 0, y, width+x, height, color)        // part that is in the screen
		return
	} else if x+width > screenWidth { // right part is off-screen
		drawRect(dst, 0, y, x+width-screenWidth, height, color) // teleported right part
		drawRect(dst, x, y, screenWidth-x, height, color)       // part that is in the screen
		return
	}

	if y < 0 { // upper part is off-screen
		drawRect(dst, x, screenHeight+y, width, -y, color) // teleported upper part
		drawRect(dst, x, 0, width, height+y, color)        // part that is in the screen
		return
	} else if y+height > screenHeight { // bottom part is off-screen
		drawRect(dst, x, 0, width, y+height-screenHeight, color) // teleported bottom part
		drawRect(dst, x, y, width, screenHeight-y, color)        // part that is in the screen
		return
	}

	// Draw the rectangle at last
	// --------------------------
	ebitenutil.DrawRect(dst, x, y, width, height, color)
}
