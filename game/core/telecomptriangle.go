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

package core

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	indices [24]uint16
)

func init() {
	for iRect := uint16(0); iRect < 4; iRect++ {
		indices[iRect*6] = iRect * 4
		indices[iRect*6+1] = iRect*4 + 2
		indices[iRect*6+2] = iRect*4 + 1
		indices[iRect*6+3] = iRect*4 + 1
		indices[iRect*6+4] = iRect*4 + 2
		indices[iRect*6+5] = iRect*4 + 3
	}

}

// TeleCompTriang is TeleComp with triangulation information for the DrawTriangles and DrawTriangleShader methods
type TeleCompTriang struct {
	TeleComp
	vertices [16]ebiten.Vertex
	color    [4]float32
}

func (t *TeleCompTriang) SetColor(clr *color.RGBA) {
	t.color = [4]float32{float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0, float32(clr.A) / 255.0}
}

func (t *TeleCompTriang) Update(pureRect *RectF32) {
	t.TeleComp.Update(pureRect)
	t.updateVertices()
}

func (t *TeleCompTriang) updateVertices() {
	var offset uint16
	for iRect := uint8(0); iRect < t.NumRects; iRect++ {
		rect := &t.Rects[iRect]

		rightX := rect.Pos.X + rect.Size.X
		rightXInUnit := rect.PosInUnit.X + rect.Size.X

		bottomY := rect.Pos.Y + rect.Size.Y
		bottomYInUnit := rect.PosInUnit.Y + rect.Size.Y

		t.vertices[offset] = ebiten.Vertex{ // Top Left corner
			DstX:   rect.Pos.X,
			DstY:   rect.Pos.Y,
			SrcX:   rect.PosInUnit.X,
			SrcY:   rect.PosInUnit.Y,
			ColorR: t.color[0],
			ColorG: t.color[1],
			ColorB: t.color[2],
			ColorA: t.color[3],
		}
		t.vertices[offset+1] = ebiten.Vertex{ // Top Right Corner
			DstX:   rightX,
			DstY:   rect.Pos.Y,
			SrcX:   rightXInUnit,
			SrcY:   rect.PosInUnit.Y,
			ColorR: t.color[0],
			ColorG: t.color[1],
			ColorB: t.color[2],
			ColorA: t.color[3],
		}
		t.vertices[offset+2] = ebiten.Vertex{ // Bottom Left Corner
			DstX:   rect.Pos.X,
			DstY:   bottomY,
			SrcX:   rect.PosInUnit.X,
			SrcY:   bottomYInUnit,
			ColorR: t.color[0],
			ColorG: t.color[1],
			ColorB: t.color[2],
			ColorA: t.color[3],
		}
		t.vertices[offset+3] = ebiten.Vertex{ // Bottom Right Corner
			DstX:   rightX,
			DstY:   bottomY,
			SrcX:   rightXInUnit,
			SrcY:   bottomYInUnit,
			ColorR: t.color[0],
			ColorG: t.color[1],
			ColorB: t.color[2],
			ColorA: t.color[3],
		}

		offset += 4
	}
}

func (t *TeleCompTriang) Triangles() ([]ebiten.Vertex, []uint16) {
	return t.vertices[:t.NumRects*4], indices[:t.NumRects*6]
}
