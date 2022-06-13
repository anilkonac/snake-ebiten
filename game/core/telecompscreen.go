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

//  Teleport component to be drawn on the screen.
type TeleCompScreen struct {
	TeleComp
	vertices [16]ebiten.Vertex
	color    [4]float32
}

func (t *TeleCompScreen) SetColor(clr *color.RGBA) {
	t.color = [4]float32{float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0, float32(clr.A) / 255.0}
}

func (t *TeleCompScreen) Update(pureRect *RectF32) {
	t.TeleComp.Update(pureRect)
	t.updateVertices()
}

func (t *TeleCompScreen) updateVertices() {
	var offset uint16
	for iRect := uint8(0); iRect < t.NumRects; iRect++ {
		rect := &t.Rects[iRect]
		rightXDest := rect.Pos.X + rect.Size.X
		rightXSrc := rect.PosInUnit.X + rect.Size.X
		bottomYDest := rect.Pos.Y + rect.Size.Y
		bottomYSrc := rect.PosInUnit.Y + rect.Size.Y

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
			DstX:   rightXDest,
			DstY:   rect.Pos.Y,
			SrcX:   rightXSrc,
			SrcY:   rect.PosInUnit.Y,
			ColorR: t.color[0],
			ColorG: t.color[1],
			ColorB: t.color[2],
			ColorA: t.color[3],
		}
		t.vertices[offset+2] = ebiten.Vertex{ // Bottom Left Corner
			DstX:   rect.Pos.X,
			DstY:   bottomYDest,
			SrcX:   rect.PosInUnit.X,
			SrcY:   bottomYSrc,
			ColorR: t.color[0],
			ColorG: t.color[1],
			ColorB: t.color[2],
			ColorA: t.color[3],
		}
		t.vertices[offset+3] = ebiten.Vertex{ // Bottom Right Corner
			DstX:   rightXDest,
			DstY:   bottomYDest,
			SrcX:   rightXSrc,
			SrcY:   bottomYSrc,
			ColorR: t.color[0],
			ColorG: t.color[1],
			ColorB: t.color[2],
			ColorA: t.color[3],
		}

		offset += 4
	}
}

func (t *TeleCompScreen) Triangles() ([]ebiten.Vertex, []uint16) {
	return t.vertices[:t.NumRects*4], indices[:t.NumRects*6]
}
