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

package snake

import (
	"image/color"

	"github.com/anilkonac/snake-ebiten/game/param"
	t "github.com/anilkonac/snake-ebiten/game/tool"
	"github.com/hajimehoshi/ebiten/v2"
)

// Teleportable/Teleport unit
type teleUnit struct {
	rects    [4]t.RectF32
	numRects uint8
}

func newTeleUnit(pureRect *t.RectF32) *teleUnit {
	tUnit := new(teleUnit)
	tUnit.split(pureRect)
	return tUnit
}

func (t teleUnit) init(pureRect *t.RectF32) {
	t.numRects = 0
	t.split(pureRect)
}

func (tu *teleUnit) split(rect *t.RectF32) {
	if (rect.Size.X <= 0) || (rect.Size.Y <= 0) {
		return
	}

	if !param.TeleportActive {
		tu.rects[tu.numRects] = *rect
		tu.numRects++
		return
	}

	rightX := rect.Pos.X + rect.Size.X
	bottomY := rect.Pos.Y + rect.Size.Y

	if rect.Pos.X < 0 { // left part is off-screen
		tu.split(
			&t.RectF32{ // teleported left part
				Pos:       t.Vec32{rect.Pos.X + param.ScreenWidth, rect.Pos.Y},
				Size:      t.Vec32{-rect.Pos.X, rect.Size.Y},
				PosInUnit: t.Vec32{0, 0},
			})

		tu.split(
			&t.RectF32{ // part in the screen
				Pos:       t.Vec32{0, rect.Pos.Y},
				Size:      t.Vec32{rightX, rect.Size.Y},
				PosInUnit: t.Vec32{-rect.Pos.X, 0},
			})

		return
	} else if rightX > param.ScreenWidth { // right part is off-screen
		tu.split(
			&t.RectF32{ // teleported right part
				Pos:       t.Vec32{0, rect.Pos.Y},
				Size:      t.Vec32{rightX - param.ScreenWidth, rect.Size.Y},
				PosInUnit: t.Vec32{param.ScreenWidth - rect.Pos.X, 0},
			})

		tu.split(
			&t.RectF32{ // part in the screen
				Pos:       t.Vec32{rect.Pos.X, rect.Pos.Y},
				Size:      t.Vec32{param.ScreenWidth - rect.Pos.X, rect.Size.Y},
				PosInUnit: t.Vec32{0, 0},
			})

		return
	}

	if rect.Pos.Y < 0 { // upper part is off-screen
		tu.split(
			&t.RectF32{ // teleported upper part
				Pos:       t.Vec32{rect.Pos.X, param.ScreenHeight + rect.Pos.Y},
				Size:      t.Vec32{rect.Size.X, -rect.Pos.Y},
				PosInUnit: t.Vec32{rect.PosInUnit.X, 0},
			})

		tu.split(
			&t.RectF32{ // part in the screen
				Pos:       t.Vec32{rect.Pos.X, 0},
				Size:      t.Vec32{rect.Size.X, bottomY},
				PosInUnit: t.Vec32{rect.PosInUnit.X, -rect.Pos.Y},
			})

		return
	} else if bottomY > param.ScreenHeight { // bottom part is off-screen
		tu.split(
			&t.RectF32{ // teleported bottom part
				Pos:       t.Vec32{rect.Pos.X, 0},
				Size:      t.Vec32{rect.Size.X, bottomY - param.ScreenHeight},
				PosInUnit: t.Vec32{rect.PosInUnit.X, param.ScreenHeight - rect.Pos.Y},
			})

		tu.split(
			&t.RectF32{ // part in the screen
				Pos:       t.Vec32{rect.Pos.X, rect.Pos.Y},
				Size:      t.Vec32{rect.Size.X, param.ScreenHeight - rect.Pos.Y},
				PosInUnit: t.Vec32{rect.PosInUnit.X, 0},
			})

		return
	}

	// Add the split rectangle to the rects array
	tu.rects[tu.numRects] = *rect
	tu.numRects++
}

//  Teleport unit to be drawn on the screen.
type teleUnitScreen struct {
	teleUnit
	vertices []ebiten.Vertex
	indices  []uint16
}

func newTeleDrawUnit(pureRect *t.RectF32) *teleUnitScreen {
	unit := new(teleUnitScreen)
	unit.update(pureRect)
	return unit
}

func (t *teleUnitScreen) update(pureRect *t.RectF32) {
	t.teleUnit.init(pureRect)
	// TODO: Update vertices and indices
}

func (t *teleUnitScreen) triangles(clr *color.RGBA) (vertices []ebiten.Vertex, indices []uint16) {
	return t.vertices, t.indices
}
