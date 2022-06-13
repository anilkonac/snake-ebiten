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
	Rects    [4]RectF32
	NumRects uint8
}

func (t *TeleComp) Update(pureRect *RectF32) {
	t.NumRects = 0
	t.split(pureRect)
}

// Divide rectangle up to 4 based on where it is off-screen.
func (t *TeleComp) split(rect *RectF32) {
	if (rect.Size.X <= 0) || (rect.Size.Y <= 0) {
		return
	}

	if !param.TeleportActive {
		t.Rects[0] = *rect
		t.NumRects = 1
		return
	}

	rightX := rect.Pos.X + rect.Size.X
	bottomY := rect.Pos.Y + rect.Size.Y

	slicedLeft := (rect.Pos.X < 0)
	slicedRight := (rightX > param.ScreenWidth)
	slicedUp := (rect.Pos.Y < 0)
	slicedDown := (bottomY > param.ScreenHeight)

	switch {
	case slicedUp && slicedLeft:
		t.splitTopLeft(rect)
	case slicedUp && slicedRight:
		t.splitTopRight(rect)
	case slicedDown && slicedLeft:
		t.splitBottomLeft(rect)
	case slicedDown && slicedRight:
		t.splitBottomRight(rect)
	case slicedUp:
		t.splitUp(rect)
	case slicedDown:
		t.splitDown(rect)
	case slicedLeft:
		t.splitLeft(rect)
	case slicedRight:
		t.splitRight(rect)
	default:
		t.Rects[0] = *rect
		t.NumRects = 1
	}
}

func (t *TeleComp) splitTopLeft(rect *RectF32) {
	rightRectsWidth := rect.Size.X + rect.Pos.X
	upperRectsPosY := param.ScreenHeight + rect.Pos.Y
	leftRectsPosX := param.ScreenWidth + rect.Pos.X
	lowerRectsHeight := rect.Size.Y + rect.Pos.Y

	t.Rects[0] = RectF32{ // teleported upper left part
		Pos:       Vec32{leftRectsPosX, upperRectsPosY},
		Size:      Vec32{-rect.Pos.X, -rect.Pos.Y},
		PosInUnit: Vec32{0, 0},
	}

	t.Rects[1] = RectF32{ // teleported upper right part
		Pos:       Vec32{0, upperRectsPosY},
		Size:      Vec32{rightRectsWidth, -rect.Pos.Y},
		PosInUnit: Vec32{-rect.Pos.X, 0},
	}

	t.Rects[2] = RectF32{ // teleported lower left part
		Pos:       Vec32{leftRectsPosX, 0},
		Size:      Vec32{-rect.Pos.X, lowerRectsHeight},
		PosInUnit: Vec32{0, -rect.Pos.Y},
	}

	t.Rects[3] = RectF32{ // lower right part that in the screen
		Pos:       Vec32{0, 0},
		Size:      Vec32{rightRectsWidth, lowerRectsHeight},
		PosInUnit: Vec32{-rect.Pos.X, -rect.Pos.Y},
	}

	t.NumRects = 4
}

func (t *TeleComp) splitTopRight(rect *RectF32) {
	leftRectsWidth := param.ScreenWidth - rect.Pos.X
	upperRectsY := param.ScreenHeight + rect.Pos.Y
	rightRectsWidth := rect.Size.X - leftRectsWidth
	lowerY := rect.Size.Y + rect.Pos.Y

	t.Rects[0] = RectF32{ // teleported upper left part
		Pos:       Vec32{rect.Pos.X, upperRectsY},
		Size:      Vec32{leftRectsWidth, -rect.Pos.Y},
		PosInUnit: Vec32{0, 0},
	}
	t.Rects[1] = RectF32{ // teleported upper right part
		Pos:       Vec32{0, upperRectsY},
		Size:      Vec32{rightRectsWidth, -rect.Pos.Y},
		PosInUnit: Vec32{leftRectsWidth, 0},
	}
	t.Rects[2] = RectF32{ // lower left part that in the screen
		Pos:       Vec32{rect.Pos.X, 0},
		Size:      Vec32{leftRectsWidth, lowerY},
		PosInUnit: Vec32{0, -rect.Pos.Y},
	}
	t.Rects[3] = RectF32{ // teleported lower right part
		Pos:       Vec32{0, 0},
		Size:      Vec32{rightRectsWidth, lowerY},
		PosInUnit: Vec32{leftRectsWidth, -rect.Pos.Y},
	}

	t.NumRects = 4
}

func (t *TeleComp) splitBottomLeft(rect *RectF32) {
	leftRectsX := param.ScreenWidth + rect.Pos.X
	upperRectsHeight := param.ScreenHeight - rect.Pos.Y
	lowerRectsHeight := rect.Size.Y - upperRectsHeight
	rightRectsWidth := rect.Size.X + rect.Pos.X

	t.Rects[0] = RectF32{ // teleported top left part
		Pos:       Vec32{leftRectsX, rect.Pos.Y},
		Size:      Vec32{-rect.Pos.X, upperRectsHeight},
		PosInUnit: Vec32{0, 0},
	}
	t.Rects[1] = RectF32{ // top right part in the screen
		Pos:       Vec32{0, rect.Pos.Y},
		Size:      Vec32{rightRectsWidth, upperRectsHeight},
		PosInUnit: Vec32{-rect.Pos.X, 0},
	}
	t.Rects[2] = RectF32{ // teleported bottom left part
		Pos:       Vec32{leftRectsX, 0},
		Size:      Vec32{-rect.Pos.X, lowerRectsHeight},
		PosInUnit: Vec32{0, upperRectsHeight},
	}
	t.Rects[3] = RectF32{ // teleported bottom right part
		Pos:       Vec32{0, 0},
		Size:      Vec32{rightRectsWidth, lowerRectsHeight},
		PosInUnit: Vec32{-rect.Pos.X, upperRectsHeight},
	}

	t.NumRects = 4
}

func (t *TeleComp) splitBottomRight(rect *RectF32) {
	leftRectsWidth := param.ScreenWidth - rect.Pos.X
	rightRectsWidth := rect.Size.X - leftRectsWidth
	upperRectsHeight := param.ScreenHeight - rect.Pos.Y
	lowerRectsHeight := rect.Size.Y - upperRectsHeight

	t.Rects[0] = RectF32{ // top left part in the screen
		Pos:       Vec32{rect.Pos.X, rect.Pos.Y},
		Size:      Vec32{leftRectsWidth, upperRectsHeight},
		PosInUnit: Vec32{0, 0},
	}
	t.Rects[1] = RectF32{ // teleported top right part
		Pos:       Vec32{0, rect.Pos.Y},
		Size:      Vec32{rightRectsWidth, upperRectsHeight},
		PosInUnit: Vec32{leftRectsWidth, 0},
	}
	t.Rects[2] = RectF32{ // teleported bottom left part
		Pos:       Vec32{rect.Pos.X, 0},
		Size:      Vec32{leftRectsWidth, lowerRectsHeight},
		PosInUnit: Vec32{0, upperRectsHeight},
	}
	t.Rects[3] = RectF32{ // teleported bottom right part
		Pos:       Vec32{0, 0},
		Size:      Vec32{rightRectsWidth, lowerRectsHeight},
		PosInUnit: Vec32{leftRectsWidth, upperRectsHeight},
	}

	t.NumRects = 4
}

func (t *TeleComp) splitUp(rect *RectF32) {
	bottomY := rect.Pos.Y + rect.Size.Y

	t.Rects[0] = RectF32{ // teleported upper part
		Pos:       Vec32{rect.Pos.X, param.ScreenHeight + rect.Pos.Y},
		Size:      Vec32{rect.Size.X, -rect.Pos.Y},
		PosInUnit: Vec32{0, 0},
	}
	t.Rects[1] = RectF32{ // part in the screen
		Pos:       Vec32{rect.Pos.X, 0},
		Size:      Vec32{rect.Size.X, bottomY},
		PosInUnit: Vec32{0, -rect.Pos.Y},
	}

	t.NumRects = 2
}

func (t *TeleComp) splitDown(rect *RectF32) {
	bottomY := rect.Pos.Y + rect.Size.Y

	t.Rects[0] = RectF32{ // part in the screen
		Pos:       Vec32{rect.Pos.X, rect.Pos.Y},
		Size:      Vec32{rect.Size.X, param.ScreenHeight - rect.Pos.Y},
		PosInUnit: Vec32{0, 0},
	}
	t.Rects[1] = RectF32{ // teleported bottom part
		Pos:       Vec32{rect.Pos.X, 0},
		Size:      Vec32{rect.Size.X, bottomY - param.ScreenHeight},
		PosInUnit: Vec32{0, param.ScreenHeight - rect.Pos.Y},
	}

	t.NumRects = 2
}

func (t *TeleComp) splitLeft(rect *RectF32) {
	rightX := rect.Pos.X + rect.Size.X

	t.Rects[0] = RectF32{ // teleported left part
		Pos:       Vec32{rect.Pos.X + param.ScreenWidth, rect.Pos.Y},
		Size:      Vec32{-rect.Pos.X, rect.Size.Y},
		PosInUnit: Vec32{0, 0},
	}
	t.Rects[1] = RectF32{ // part in the screen
		Pos:       Vec32{0, rect.Pos.Y},
		Size:      Vec32{rightX, rect.Size.Y},
		PosInUnit: Vec32{-rect.Pos.X, 0},
	}

	t.NumRects = 2
}

func (t *TeleComp) splitRight(rect *RectF32) {
	rightX := rect.Pos.X + rect.Size.X

	t.Rects[0] = RectF32{ // part in the screen
		Pos:       Vec32{rect.Pos.X, rect.Pos.Y},
		Size:      Vec32{param.ScreenWidth - rect.Pos.X, rect.Size.Y},
		PosInUnit: Vec32{0, 0},
	}

	t.Rects[1] = RectF32{ // teleported right part
		Pos:       Vec32{0, rect.Pos.Y},
		Size:      Vec32{rightX - param.ScreenWidth, rect.Size.Y},
		PosInUnit: Vec32{param.ScreenWidth - rect.Pos.X, 0},
	}

	t.NumRects = 2
}
