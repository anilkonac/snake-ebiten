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

// Rectangle compatible with float32 type parameters of the ebiten.DrawTriangleShader function.
type rectF32 struct {
	pos       vec32
	size      vec32
	posInUnit vec32
}

func newRect(pos, size vec32) *rectF32 {
	return &rectF32{pos, size, vec32{0, 0}}
}

// Divide rectangle up to 4 based on where it is off-screen.
func (r rectF32) split(rects *[]rectF32) {
	if (r.size.x <= 0) || (r.size.y <= 0) {
		return
	}

	rightX := r.pos.x + r.size.x
	bottomY := r.pos.y + r.size.y

	if r.pos.x < 0 { // left part is off-screen
		rectF32{vec32{r.pos.x + ScreenWidth, r.pos.y}, vec32{-r.pos.x, r.size.y}, vec32{0, 0}}.split(rects) // teleported left part
		rectF32{vec32{0, r.pos.y}, vec32{rightX, r.size.y}, vec32{-r.pos.x, 0}}.split(rects)                // part in the screen
		return
	} else if rightX > ScreenWidth { // right part is off-screen
		rectF32{vec32{0, r.pos.y}, vec32{rightX - ScreenWidth, r.size.y}, vec32{ScreenWidth - r.pos.x, 0}}.split(rects) // teleported right part
		rectF32{vec32{r.pos.x, r.pos.y}, vec32{ScreenWidth - r.pos.x, r.size.y}, vec32{0, 0}}.split(rects)              // part in the screen
		return
	}

	if r.pos.y < 0 { // upper part is off-screen
		rectF32{vec32{r.pos.x, ScreenHeight + r.pos.y}, vec32{r.size.x, -r.pos.y}, vec32{r.posInUnit.x, 0}}.split(rects) // teleported upper part
		rectF32{vec32{r.pos.x, 0}, vec32{r.size.x, bottomY}, vec32{r.posInUnit.x, -r.pos.y}}.split(rects)                // part in the screen
		return
	} else if bottomY > ScreenHeight { // bottom part is off-screen
		rectF32{vec32{r.pos.x, 0}, vec32{r.size.x, bottomY - ScreenHeight}, vec32{r.posInUnit.x, ScreenHeight - r.pos.y}}.split(rects) // teleported bottom part
		rectF32{vec32{r.pos.x, r.pos.y}, vec32{r.size.x, ScreenHeight - r.pos.y}, vec32{r.posInUnit.x, 0}}.split(rects)                // part in the screen
		return
	}

	// Add the split rectangle to the rects slice.
	*rects = append(*rects, r)
}

func intersects(rectA, rectB *rectF32, tolerance float32) bool {
	aRightX := rectA.pos.x + rectA.size.x
	bRightX := rectB.pos.x + rectB.size.x
	aBottomY := rectA.pos.y + rectA.size.y
	bBottomY := rectB.pos.y + rectB.size.y

	if (rectA.pos.x-rectB.pos.x <= tolerance) && (aRightX-rectB.pos.x <= tolerance) { // rectA is on the left side of rectB
		return false
	}

	if (rectA.pos.x-bRightX >= -tolerance) && (aRightX-bRightX >= -tolerance) { // rectA is on the right side of rectB
		return false
	}

	if (rectA.pos.y-rectB.pos.y <= tolerance) && (aBottomY-rectB.pos.y <= tolerance) { // rectA is above rectB
		return false
	}

	if (rectA.pos.y-bBottomY >= -tolerance) && (aBottomY-bBottomY >= -tolerance) { // rectA is under rectB
		return false
	}

	return true
}

func (r rectF32) vertices(clr *color.RGBA) []ebiten.Vertex {
	var fR, fG, fB, fA float32 = float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0, float32(clr.A) / 255.0
	return []ebiten.Vertex{
		{ // Top Left corner
			DstX:   r.pos.x,
			DstY:   r.pos.y,
			SrcX:   r.posInUnit.x,
			SrcY:   r.posInUnit.y,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
		{ // Top Right Corner
			DstX:   r.pos.x + r.size.x,
			DstY:   r.pos.y,
			SrcX:   r.posInUnit.x + r.size.x,
			SrcY:   r.posInUnit.y,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
		{ // Bottom Left Corner
			DstX:   r.pos.x,
			DstY:   r.pos.y + r.size.y,
			SrcX:   r.posInUnit.x,
			SrcY:   r.posInUnit.y + r.size.y,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
		{ // Bottom Right Corner
			DstX:   r.pos.x + r.size.x,
			DstY:   r.pos.y + r.size.y,
			SrcX:   r.posInUnit.x + r.size.x,
			SrcY:   r.posInUnit.y + r.size.y,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
	}
}

func (r rectF32) drawOuterRect(dst *ebiten.Image, clr color.Color) {
	pos64 := r.pos.to64()
	size64 := r.size.to64()
	ebitenutil.DrawRect(dst, pos64.x, pos64.y, size64.x, size64.y, color.RGBA{255, 255, 255, 96})
}

func markPoint(dst *ebiten.Image, p *vec64, length float64, clr color.Color) {
	ebitenutil.DrawLine(dst, p.x-length, p.y, p.x+length, p.y, clr)
	ebitenutil.DrawLine(dst, p.x, p.y-length, p.x, p.y+length, clr)
}
