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
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type rect struct {
	x, y             int16
	width, height    int16
	xInUnit, yInUnit int16
}

// Divide rectangle up to 4 based on where it is off-screen.
func (r rect) split(rects *[]rect) {
	if r.width <= 0 || r.height <= 0 {
		return
	}

	rightX := r.x + r.width
	bottomY := r.y + r.height

	if r.x < 0 { // left part is off-screen
		rect{r.x + ScreenWidth, r.y, -r.x, r.height, 0, 0}.split(rects) // teleported left part
		rect{0, r.y, rightX, r.height, -r.x, 0}.split(rects)            // part in the screen
		return
	} else if rightX > ScreenWidth { // right part is off-screen
		rect{0, r.y, rightX - ScreenWidth, r.height, ScreenWidth - r.x, 0}.split(rects) // teleported right part
		rect{r.x, r.y, ScreenWidth - r.x, r.height, 0, 0}.split(rects)                  // part in the screen
		return
	}

	if r.y < 0 { // upper part is off-screen
		rect{r.x, ScreenHeight + r.y, r.width, -r.y, r.xInUnit, 0}.split(rects) // teleported upper part
		rect{r.x, 0, r.width, bottomY, r.xInUnit, -r.y}.split(rects)            // part in the screen
		return
	} else if bottomY > ScreenHeight { // bottom part is off-screen
		rect{r.x, 0, r.width, bottomY - ScreenHeight, r.xInUnit, ScreenHeight - r.y}.split(rects) // teleported bottom part
		rect{r.x, r.y, r.width, ScreenHeight - r.y, r.xInUnit, 0}.split(rects)                    // part in the screen
		return
	}

	// Add the split rectangle to the rects slice.
	*rects = append(*rects, r)
}

func intersects(rectA, rectB *rect, tolerance int16) bool {
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

func (r rect) vertices(color color.Color) []ebiten.Vertex {
	uR, uG, uB, uA := color.RGBA()
	fR, fG, fB, fA := float32(uR), float32(uG), float32(uB), float32(uA)
	// fX, fY := float32(r.x), float32(fY)
	vertices := []ebiten.Vertex{
		{ // Top Left corner
			DstX:   float32(r.x),
			DstY:   float32(r.y),
			SrcX:   float32(r.xInUnit),
			SrcY:   float32(r.yInUnit),
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
		{ // Top Right Corner
			DstX:   float32(r.x + r.width),
			DstY:   float32(r.y),
			SrcX:   float32(r.xInUnit + r.width),
			SrcY:   float32(r.yInUnit),
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
		{ // Bottom Left Corner
			DstX:   float32(r.x),
			DstY:   float32(r.y + r.height),
			SrcX:   float32(r.xInUnit),
			SrcY:   float32(r.yInUnit + r.height),
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
		{ // Bottom Right Corner
			DstX:   float32(r.x + r.width),
			DstY:   float32(r.y + r.height),
			SrcX:   float32(r.xInUnit + r.width),
			SrcY:   float32(r.yInUnit + r.height),
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
	}
	return vertices
}

func (r rect) drawOuterRect(dst *ebiten.Image, clr color.Color) {
	x64 := float64(r.x)
	y64 := float64(r.y)
	xRight64 := float64(r.x + r.width)
	yBottom64 := float64(r.y + r.height)
	ebitenutil.DrawLine(dst, x64, y64, xRight64, y64, clr)
	ebitenutil.DrawLine(dst, x64, yBottom64, xRight64, yBottom64, clr)
	ebitenutil.DrawLine(dst, x64+1, y64, x64+1, yBottom64, clr)
	ebitenutil.DrawLine(dst, xRight64, y64, xRight64, yBottom64, clr)
}

func markPoint(dst *ebiten.Image, pX, pY float64, clr color.Color) {
	ebitenutil.DrawLine(dst, pX-3, pY, pX+3, pY, clr)
	ebitenutil.DrawLine(dst, pX, pY-3, pX, pY+3, clr)
}
