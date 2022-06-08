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

	"github.com/anilkonac/snake-ebiten/game/param"
	t "github.com/anilkonac/snake-ebiten/game/tool"
	"github.com/hajimehoshi/ebiten/v2"
)

var indices [24]uint16

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

// Teleportable/Teleport component
type TeleComp struct {
	Rects    [4]t.RectF32
	NumRects uint8
}

func (t *TeleComp) Update(pureRect *t.RectF32) {
	t.NumRects = 0
	t.split(pureRect)
}

func (tu *TeleComp) split(rect *t.RectF32) {
	if (rect.Size.X <= 0) || (rect.Size.Y <= 0) {
		return
	}

	if !param.TeleportActive {
		tu.Rects[tu.NumRects] = *rect
		tu.NumRects++
		return
	}

	rightX := rect.Pos.X + rect.Size.X
	bottomY := rect.Pos.Y + rect.Size.Y

	if rect.Pos.X < 0 { // left part is off-screen
		tu.split(
			&t.RectF32{ // teleported left part
				Pos:       t.Vec32{X: rect.Pos.X + param.ScreenWidth, Y: rect.Pos.Y},
				Size:      t.Vec32{X: -rect.Pos.X, Y: rect.Size.Y},
				PosInUnit: t.Vec32{X: 0, Y: 0},
			})

		tu.split(
			&t.RectF32{ // part in the screen
				Pos:       t.Vec32{X: 0, Y: rect.Pos.Y},
				Size:      t.Vec32{X: rightX, Y: rect.Size.Y},
				PosInUnit: t.Vec32{X: -rect.Pos.X, Y: 0},
			})

		return
	} else if rightX > param.ScreenWidth { // right part is off-screen
		tu.split(
			&t.RectF32{ // teleported right part
				Pos:       t.Vec32{X: 0, Y: rect.Pos.Y},
				Size:      t.Vec32{X: rightX - param.ScreenWidth, Y: rect.Size.Y},
				PosInUnit: t.Vec32{X: param.ScreenWidth - rect.Pos.X, Y: 0},
			})

		tu.split(
			&t.RectF32{ // part in the screen
				Pos:       t.Vec32{X: rect.Pos.X, Y: rect.Pos.Y},
				Size:      t.Vec32{X: param.ScreenWidth - rect.Pos.X, Y: rect.Size.Y},
				PosInUnit: t.Vec32{X: 0, Y: 0},
			})

		return
	}

	if rect.Pos.Y < 0 { // upper part is off-screen
		tu.split(
			&t.RectF32{ // teleported upper part
				Pos:       t.Vec32{X: rect.Pos.X, Y: param.ScreenHeight + rect.Pos.Y},
				Size:      t.Vec32{X: rect.Size.X, Y: -rect.Pos.Y},
				PosInUnit: t.Vec32{X: rect.PosInUnit.X, Y: 0},
			})

		tu.split(
			&t.RectF32{ // part in the screen
				Pos:       t.Vec32{X: rect.Pos.X, Y: 0},
				Size:      t.Vec32{X: rect.Size.X, Y: bottomY},
				PosInUnit: t.Vec32{X: rect.PosInUnit.X, Y: -rect.Pos.Y},
			})

		return
	} else if bottomY > param.ScreenHeight { // bottom part is off-screen
		tu.split(
			&t.RectF32{ // teleported bottom part
				Pos:       t.Vec32{X: rect.Pos.X, Y: 0},
				Size:      t.Vec32{X: rect.Size.X, Y: bottomY - param.ScreenHeight},
				PosInUnit: t.Vec32{X: rect.PosInUnit.X, Y: param.ScreenHeight - rect.Pos.Y},
			})

		tu.split(
			&t.RectF32{ // part in the screen
				Pos:       t.Vec32{X: rect.Pos.X, Y: rect.Pos.Y},
				Size:      t.Vec32{X: rect.Size.X, Y: param.ScreenHeight - rect.Pos.Y},
				PosInUnit: t.Vec32{X: rect.PosInUnit.X, Y: 0},
			})

		return
	}

	// Add the split rectangle to the rects array
	tu.Rects[tu.NumRects] = *rect
	tu.NumRects++
}

//  Teleport component to be drawn on the screen.
type TeleCompScreen struct {
	TeleComp
	vertices [16]ebiten.Vertex
	color    [4]float32
}

func (t *TeleCompScreen) SetColor(clr *color.RGBA) {
	t.color = [4]float32{float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0, float32(clr.A) / 255.0}
}

func (t *TeleCompScreen) Update(pureRect *t.RectF32) {
	t.TeleComp.Update(pureRect)
	t.updateVertices()
}

func (t *TeleCompScreen) updateVertices() {
	var offset uint16
	for iRect := uint8(0); iRect < t.NumRects; iRect++ {
		rect := &t.Rects[iRect]
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
			DstX:   rect.Pos.X + rect.Size.X,
			DstY:   rect.Pos.Y,
			SrcX:   rect.PosInUnit.X + rect.Size.X,
			SrcY:   rect.PosInUnit.Y,
			ColorR: t.color[0],
			ColorG: t.color[1],
			ColorB: t.color[2],
			ColorA: t.color[3],
		}
		t.vertices[offset+2] = ebiten.Vertex{ // Bottom Left Corner
			DstX:   rect.Pos.X,
			DstY:   rect.Pos.Y + rect.Size.Y,
			SrcX:   rect.PosInUnit.X,
			SrcY:   rect.PosInUnit.Y + rect.Size.Y,
			ColorR: t.color[0],
			ColorG: t.color[1],
			ColorB: t.color[2],
			ColorA: t.color[3],
		}
		t.vertices[offset+3] = ebiten.Vertex{ // Bottom Right Corner
			DstX:   rect.Pos.X + rect.Size.X,
			DstY:   rect.Pos.Y + rect.Size.Y,
			SrcX:   rect.PosInUnit.X + rect.Size.X,
			SrcY:   rect.PosInUnit.Y + rect.Size.Y,
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
