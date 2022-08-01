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
	"github.com/anilkonac/snake-ebiten/game/param"
)

// Teleportable/Teleport component
type TeleComp struct {
	Rects    [4]RectF32
	NumRects uint8
}

func (t *TeleComp) Update(pureRect *RectF32) {
	t.NumRects = 0
	t.split(*pureRect)
}

func (t *TeleComp) split(rect RectF32) {
	if (rect.Size.X <= 0) || (rect.Size.Y <= 0) {
		return
	}

	if !param.TeleportEnabled {
		t.Rects[0] = rect
		t.NumRects = 1
		return
	}

	rightX := rect.Pos.X + rect.Size.X
	bottomY := rect.Pos.Y + rect.Size.Y

	if rect.Pos.X < 0 { // left part is off-screen
		t.split(
			RectF32{ // teleported left part
				Pos:       Vec32{rect.Pos.X + param.ScreenWidth, rect.Pos.Y},
				Size:      Vec32{-rect.Pos.X, rect.Size.Y},
				PosInUnit: Vec32{0, 0},
			})

		t.split(
			RectF32{ // part in the screen
				Pos:       Vec32{0, rect.Pos.Y},
				Size:      Vec32{rightX, rect.Size.Y},
				PosInUnit: Vec32{-rect.Pos.X, 0},
			})

		return
	} else if rightX > param.ScreenWidth { // right part is off-screen
		t.split(
			RectF32{ // teleported right part
				Pos:       Vec32{0, rect.Pos.Y},
				Size:      Vec32{rightX - param.ScreenWidth, rect.Size.Y},
				PosInUnit: Vec32{param.ScreenWidth - rect.Pos.X, 0},
			})

		t.split(
			RectF32{ // part in the screen
				Pos:       Vec32{rect.Pos.X, rect.Pos.Y},
				Size:      Vec32{param.ScreenWidth - rect.Pos.X, rect.Size.Y},
				PosInUnit: Vec32{0, 0},
			})

		return
	}

	if rect.Pos.Y < 0 { // upper part is off-screen
		t.split(
			RectF32{ // teleported upper part
				Pos:       Vec32{rect.Pos.X, param.ScreenHeight + rect.Pos.Y},
				Size:      Vec32{rect.Size.X, -rect.Pos.Y},
				PosInUnit: Vec32{rect.PosInUnit.X, 0},
			})

		t.split(
			RectF32{ // part in the screen
				Pos:       Vec32{rect.Pos.X, 0},
				Size:      Vec32{rect.Size.X, bottomY},
				PosInUnit: Vec32{rect.PosInUnit.X, -rect.Pos.Y},
			})

		return
	} else if bottomY > param.ScreenHeight { // bottom part is off-screen
		t.split(
			RectF32{ // teleported bottom part
				Pos:       Vec32{rect.Pos.X, 0},
				Size:      Vec32{rect.Size.X, bottomY - param.ScreenHeight},
				PosInUnit: Vec32{rect.PosInUnit.X, param.ScreenHeight - rect.Pos.Y},
			})

		t.split(
			RectF32{ // part in the screen
				Pos:       Vec32{rect.Pos.X, rect.Pos.Y},
				Size:      Vec32{rect.Size.X, param.ScreenHeight - rect.Pos.Y},
				PosInUnit: Vec32{rect.PosInUnit.X, 0},
			})

		return
	}

	// Add the split rectangle to the rects array
	t.Rects[t.NumRects] = rect
	t.NumRects++
}
