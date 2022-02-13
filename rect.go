package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const epsilon = 0.001

type rectF64 struct {
	x, y          float64
	width, height float64
}

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

func (a rectF64) intersects(b *rectF64) bool {

	aRects := make([]rectF64, 0, 4)
	bRects := make([]rectF64, 0, 4)
	a.slice(&aRects)
	b.slice(&bRects)

	for iARect := 0; iARect < len(aRects); iARect++ {
		aRect := &aRects[iARect]
		for iBRect := 0; iBRect < len(bRects); iBRect++ {
			bRect := &bRects[iBRect]

			aRightX := aRect.x + aRect.width
			bRightX := bRect.x + bRect.width
			aBottomY := aRect.y + aRect.height
			bBottomY := bRect.y + bRect.height

			if (aRect.x-bRect.x <= epsilon) && (aRightX-bRect.x <= epsilon) { // a is on the left side of b
				continue
			}

			if (aRect.x-bRightX >= -epsilon) && (aRightX-bRightX >= -epsilon) { // a is on the right side of b
				continue
			}

			if (aRect.y-bRect.y <= epsilon) && (aBottomY-bRect.y <= epsilon) { // a is above b
				continue
			}

			if (aRect.y-bBottomY >= -epsilon) && (aBottomY-bBottomY >= -epsilon) { // a is under b
				continue
			}

			return true
		}
	}

	return false
}

func (r rectF64) draw(dst *ebiten.Image, clr color.Color) {
	drawRect(dst, r.x, r.y, r.width, r.height, clr)
}
