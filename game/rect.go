/*
snake-ebiten
Copyright (C) 2022 Anıl Konaç

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Rectangle compatible with float64 type parameters of the ebitenutil.DrawRect function.
type rectF64 struct {
	x, y          float64
	width, height float64
}

// Divide rectangle up to 4 based on where it is off-screen.
func (r rectF64) split(rects *[]rectF64) {
	if r.width < 0 || r.height < 0 {
		return
	}

	rightX := r.x + r.width
	bottomY := r.y + r.height

	if r.x < 0 { // left part is off-screen
		rectF64{r.x + GameWidth, r.y, -r.x, r.height}.split(rects) // teleported left part
		rectF64{0, r.y, rightX, r.height}.split(rects)             // part in the screen
		return
	} else if rightX > GameWidth { // right part is off-screen
		rectF64{0, r.y, rightX - GameWidth, r.height}.split(rects) // teleported right part
		rectF64{r.x, r.y, GameWidth - r.x, r.height}.split(rects)  // part in the screen
		return
	}

	if r.y < 0 { // upper part is off-screen
		rectF64{r.x, GameHeight + r.y, r.width, -r.y}.split(rects) // teleported upper part
		rectF64{r.x, 0, r.width, bottomY}.split(rects)             // part in the screen
		return
	} else if bottomY > GameHeight { // bottom part is off-screen
		rectF64{r.x, 0, r.width, bottomY - GameHeight}.split(rects) // teleported bottom part
		rectF64{r.x, r.y, r.width, GameHeight - r.y}.split(rects)   // part in the screen
		return
	}

	// Add the split rectangle to the rects slice.
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

func intersects(rectA, rectB rectF64, tolerance float64) bool {
	aRightX := rectA.x + rectA.width
	bRightX := rectB.x + rectB.width
	aBottomY := rectA.y + rectA.height
	bBottomY := rectB.y + rectB.height

	if (rectA.x-rectB.x <= tolerance) && (aRightX-rectB.x <= tolerance) { // rectA is on the left side of rectB
		return false
	}

	if (rectA.x-bRightX >= -tolerance) && (aRightX-bRightX >= -tolerance) { // rectA is on the right side of rectB
		return false
	}

	if (rectA.y-rectB.y <= tolerance) && (aBottomY-rectB.y <= tolerance) { // rectA is above rectB
		return false
	}

	if (rectA.y-bBottomY >= -tolerance) && (aBottomY-bBottomY >= -tolerance) { // rectA is under rectB
		return false
	}

	return true
}
