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

// Rectangle compatible with float32 type parameters of the ebiten.DrawTriangleShader function.
type rectF32 struct {
	x, y             float32
	width, height    float32
	xInUnit, yInUnit float32
}

// Divide rectangle up to 4 based on where it is off-screen.
func (r rectF32) split(rects *[]rectF32) {
	if r.width <= 0 || r.height <= 0 {
		return
	}

	rightX := r.x + r.width
	bottomY := r.y + r.height

	if r.x < 0 { // left part is off-screen
		rectF32{r.x + ScreenWidth, r.y, -r.x, r.height, 0, 0}.split(rects) // teleported left part
		rectF32{0, r.y, rightX, r.height, -r.x, 0}.split(rects)            // part in the screen
		return
	} else if rightX > ScreenWidth { // right part is off-screen
		rectF32{0, r.y, rightX - ScreenWidth, r.height, ScreenWidth - r.x, 0}.split(rects) // teleported right part
		rectF32{r.x, r.y, ScreenWidth - r.x, r.height, 0, 0}.split(rects)                  // part in the screen
		return
	}

	if r.y < 0 { // upper part is off-screen
		rectF32{r.x, ScreenHeight + r.y, r.width, -r.y, r.xInUnit, 0}.split(rects) // teleported upper part
		rectF32{r.x, 0, r.width, bottomY, r.xInUnit, -r.y}.split(rects)            // part in the screen
		return
	} else if bottomY > ScreenHeight { // bottom part is off-screen
		rectF32{r.x, 0, r.width, bottomY - ScreenHeight, r.xInUnit, ScreenHeight - r.y}.split(rects) // teleported bottom part
		rectF32{r.x, r.y, r.width, ScreenHeight - r.y, r.xInUnit, 0}.split(rects)                    // part in the screen
		return
	}

	// Add the split rectangle to the rects slice.
	*rects = append(*rects, r)
}

func (r rectF32) draw(dst *ebiten.Image, clr color.Color) {
	cr, cg, cb, ca := clr.RGBA()
	fR, fG, fB, fA := float32(cr), float32(cg), float32(cb), float32(ca)
	vertices := []ebiten.Vertex{
		{
			DstX:   r.x,
			DstY:   r.y,
			SrcX:   0,
			SrcY:   0,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
		{
			DstX:   r.x + r.width,
			DstY:   r.y,
			SrcX:   r.width,
			SrcY:   0,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
		{
			DstX:   r.x,
			DstY:   r.y + r.height,
			SrcX:   0,
			SrcY:   r.height,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
		{
			DstX:   r.x + r.width,
			DstY:   r.y + r.height,
			SrcX:   r.width,
			SrcY:   r.height,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
	}

	indices := []uint16{
		1, 0, 2,
		2, 3, 1,
	}
	op := &ebiten.DrawTrianglesShaderOptions{}

	dst.DrawTrianglesShader(vertices, indices, shaderMap[shaderBasic], op)

	if debugUnits {
		ebitenutil.DebugPrintAt(dst, fmt.Sprintf("%3.3f, %3.3f", r.x, r.y), int(r.x)-90, int(r.y)-15)
		bottomX := r.x + r.width
		bottomY := r.y + r.height
		ebitenutil.DebugPrintAt(dst, fmt.Sprintf("%3.3f, %3.3f", bottomX, bottomY), int(bottomX), int(bottomY))
	}
}

func intersects(rectA, rectB rectF32, tolerance float32) bool {
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
