package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type rectF64 struct {
	x, y          float64
	width, height float64
}

const epsilon = 0.001

// Divide rectangle up to 4 ones according to where it is off-screen.
func (r rectF64) slice(rects *[]rectF64) {
	rightX := r.x + r.width
	bottomY := r.y + r.height

	if r.x < 0 { // left part is off-screen
		rectF64{r.x + screenWidth, r.y, -r.x, r.height}.slice(rects) // teleported left part
		rectF64{0, r.y, rightX, r.height}.slice(rects)               // part in the screen
		return
	} else if rightX > screenWidth { // right part is off-screen
		rectF64{0, r.y, rightX - screenWidth, r.height}.slice(rects) // teleported right part
		rectF64{r.x, r.y, screenWidth - r.x, r.height}.slice(rects)  // part in the screen
		return
	}

	if r.y < 0 { // upper part is off-screen
		rectF64{r.x, screenHeight + r.y, r.width, -r.y}.slice(rects) // teleported upper part
		rectF64{r.x, 0, r.width, bottomY}.slice(rects)               // part in the screen
		return
	} else if bottomY > screenHeight { // bottom part is off-screen
		rectF64{r.x, 0, r.width, bottomY - screenHeight}.slice(rects) // teleported bottom part
		rectF64{r.x, r.y, r.width, screenHeight - r.y}.slice(rects)   // part in the screen
		return
	}

	// Add sliced rectangle to the slice
	*rects = append(*rects, r)
}

func (r rectF64) draw(dst *ebiten.Image, clr color.Color) {
	ebitenutil.DrawRect(dst, r.x, r.y, r.width, r.height, clr)
	if debugUnits {
		ebitenutil.DebugPrintAt(dst, fmt.Sprintf("%3.3f, %3.3f", r.x, r.y), int(r.x)-90, int(r.y)-15)
		bottomX := r.x + r.width
		bottomY := r.y + r.height
		ebitenutil.DebugPrintAt(dst, fmt.Sprintf("%3.3f, %3.3f", bottomX, bottomY), int(bottomX), int(bottomY))
	}
}
