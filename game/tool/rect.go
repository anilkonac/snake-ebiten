/*
Copyright (C) 2022 Anıl Konaç

This file is part of snake-ebiten.

snake-ebiten is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

snake-ebiten is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with snake-ebiten. If not, see <https://www.gnu.org/licenses/>.
*/

package tool

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Rectangle compatible with float32 type parameters of the ebiten.DrawTriangleShader function.
type RectF32 struct {
	Pos       Vec32
	Size      Vec32
	PosInUnit Vec32
}

func NewRect(pos, size Vec32) *RectF32 {
	return &RectF32{pos, size, Vec32{0, 0}}
}

func Intersects(rectA, rectB *RectF32, tolerance float32) bool {
	aRightX := rectA.Pos.X + rectA.Size.X
	bRightX := rectB.Pos.X + rectB.Size.X
	aBottomY := rectA.Pos.Y + rectA.Size.Y
	bBottomY := rectB.Pos.Y + rectB.Size.Y

	if (rectA.Pos.X-rectB.Pos.X <= tolerance) && (aRightX-rectB.Pos.X <= tolerance) { // rectA is on the left side of rectB
		return false
	}

	if (rectA.Pos.X-bRightX >= -tolerance) && (aRightX-bRightX >= -tolerance) { // rectA is on the right side of rectB
		return false
	}

	if (rectA.Pos.Y-rectB.Pos.Y <= tolerance) && (aBottomY-rectB.Pos.Y <= tolerance) { // rectA is above rectB
		return false
	}

	if (rectA.Pos.Y-bBottomY >= -tolerance) && (aBottomY-bBottomY >= -tolerance) { // rectA is under rectB
		return false
	}

	return true
}

func (r RectF32) Vertices(clr *color.RGBA) []ebiten.Vertex {
	fR, fG, fB, fA := float32(clr.R)/255.0, float32(clr.G)/255.0, float32(clr.B)/255.0, float32(clr.A)/255.0
	return []ebiten.Vertex{
		{ // Top Left corner
			DstX:   r.Pos.X,
			DstY:   r.Pos.Y,
			SrcX:   r.PosInUnit.X,
			SrcY:   r.PosInUnit.Y,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
		{ // Top Right Corner
			DstX:   r.Pos.X + r.Size.X,
			DstY:   r.Pos.Y,
			SrcX:   r.PosInUnit.X + r.Size.X,
			SrcY:   r.PosInUnit.Y,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
		{ // Bottom Left Corner
			DstX:   r.Pos.X,
			DstY:   r.Pos.Y + r.Size.Y,
			SrcX:   r.PosInUnit.X,
			SrcY:   r.PosInUnit.Y + r.Size.Y,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
		{ // Bottom Right Corner
			DstX:   r.Pos.X + r.Size.X,
			DstY:   r.Pos.Y + r.Size.Y,
			SrcX:   r.PosInUnit.X + r.Size.X,
			SrcY:   r.PosInUnit.Y + r.Size.Y,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		},
	}
}

func (r RectF32) DrawOuterRect(dst *ebiten.Image, clr color.Color) {
	pos64 := r.Pos.To64()
	size64 := r.Size.To64()
	ebitenutil.DrawRect(dst, pos64.X, pos64.Y, size64.X, size64.Y, color.RGBA{255, 255, 255, 96})
}

func MarkPoint(dst *ebiten.Image, p Vec64, length float64, clr color.Color) {
	ebitenutil.DrawLine(dst, p.X-length, p.Y, p.X+length, p.Y, clr)
	ebitenutil.DrawLine(dst, p.X, p.Y-length, p.X, p.Y+length, clr)
}
