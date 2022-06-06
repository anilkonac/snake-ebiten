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

var indices = [24]uint16{
	1, 0, 2,
	2, 3, 1,
	5, 4, 6,
	6, 7, 5,
	9, 8, 10,
	10, 11, 9,
	13, 12, 14,
	14, 15, 13,
}

// Teleportable/Teleport unit
type TeleUnit struct {
	Rects    [4]t.RectF32
	NumRects uint8
}

// func newTeleUnit(pureRect *t.RectF32) *TeleUnit {
// 	tUnit := new(TeleUnit)
// 	tUnit.split(pureRect)
// 	return tUnit
// }

func (t TeleUnit) Init(pureRect *t.RectF32) {
	t.NumRects = 0
	t.split(pureRect)
}

func (tu *TeleUnit) split(rect *t.RectF32) {
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

//  Teleport unit to be drawn on the screen.
type TeleUnitScreen struct {
	TeleUnit
	Vertices [16]ebiten.Vertex
	Indices  [24]uint16
}

// func newTeleDrawUnit(pureRect *t.RectF32, clr *color.RGBA) *TeleUnitScreen {
// 	unit := new(TeleUnitScreen)
// 	unit.Init(pureRect, clr)
// 	return unit
// }

func (t *TeleUnitScreen) Init(pureRect *t.RectF32, clr *color.RGBA) {
	t.TeleUnit.Init(pureRect)
	t.updateVertices(clr)
}

func (t *TeleUnitScreen) updateVertices(clr *color.RGBA) {
	fR, fG, fB, fA := float32(clr.R)/255.0, float32(clr.G)/255.0, float32(clr.B)/255.0, float32(clr.A)/255.0

	var offset uint16
	for iRect := uint8(0); iRect < t.NumRects; iRect++ {
		rect := &t.Rects[iRect]
		t.Vertices[offset] = ebiten.Vertex{ // Top Left corner
			DstX:   rect.Pos.X,
			DstY:   rect.Pos.Y,
			SrcX:   rect.PosInUnit.X,
			SrcY:   rect.PosInUnit.Y,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		}
		t.Vertices[offset+1] = ebiten.Vertex{ // Top Right Corner
			DstX:   rect.Pos.X + rect.Size.X,
			DstY:   rect.Pos.Y,
			SrcX:   rect.PosInUnit.X + rect.Size.X,
			SrcY:   rect.PosInUnit.Y,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		}
		t.Vertices[offset+2] = ebiten.Vertex{ // Bottom Left Corner
			DstX:   rect.Pos.X,
			DstY:   rect.Pos.Y + rect.Size.Y,
			SrcX:   rect.PosInUnit.X,
			SrcY:   rect.PosInUnit.Y + rect.Size.Y,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		}
		t.Vertices[offset+3] = ebiten.Vertex{ // Bottom Right Corner
			DstX:   rect.Pos.X + rect.Size.X,
			DstY:   rect.Pos.Y + rect.Size.Y,
			SrcX:   rect.PosInUnit.X + rect.Size.X,
			SrcY:   rect.PosInUnit.Y + rect.Size.Y,
			ColorR: fR,
			ColorG: fG,
			ColorB: fB,
			ColorA: fA,
		}

		offset += 4
	}
}

// func (t *TeleUnitScreen) triangles(clr *color.RGBA) (vertices []ebiten.Vertex, indices []uint16) {
// 	return t.Vertices[:], t.Indices[:t.NumRects*6]
// }
