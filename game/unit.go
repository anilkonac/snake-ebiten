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

type unit struct {
	headCenterX float32
	headCenterY float32
	length      float32
	direction   directionT
	rects       []rectF32 // rectangles that are used for both collision checking and drawing
	color       *color.RGBA
	next        *unit
	prev        *unit
}

func newUnit(headCenterX, headCenterY, length float32, direction directionT, color *color.RGBA) *unit {
	newUnit := &unit{
		headCenterX: headCenterX,
		headCenterY: headCenterY,
		length:      length,
		direction:   direction,
		color:       color,
	}
	newUnit.creteRects()

	return newUnit
}

func (u *unit) creteRects() {
	// Create the rectangle to split.
	var pureRect rectF32
	lengthF32 := float32(u.length)
	switch u.direction {
	case directionRight:
		pureRect = rectF32{
			x:      u.headCenterX - lengthF32 + halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  lengthF32,
			height: snakeWidth,
		}
	case directionLeft:
		pureRect = rectF32{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  lengthF32,
			height: snakeWidth,
		}
	case directionUp:
		pureRect = rectF32{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  snakeWidth,
			height: lengthF32,
		}
	case directionDown:
		pureRect = rectF32{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - lengthF32 + halfSnakeWidth,
			width:  snakeWidth,
			height: lengthF32,
		}
	default:
		panic("Wrong unit direction!!")
	}

	u.rects = make([]rectF32, 0, 4) // Remove old rectangles
	pureRect.split(&u.rects)        // Create split rectangles on screen edges.
}

func (u *unit) moveUp(dist float32) {
	u.headCenterY -= dist

	// teleport if head center is offscreen.
	if u.headCenterY < 0 {
		u.headCenterY += ScreenHeight
	}
}

func (u *unit) moveDown(dist float32) {
	u.headCenterY += dist

	// teleport if head center is offscreen.
	if u.headCenterY > ScreenHeight {
		u.headCenterY -= ScreenHeight
	}
}

func (u *unit) moveRight(dist float32) {
	u.headCenterX += dist

	// teleport if head center is offscreen.
	if u.headCenterX > ScreenWidth {
		u.headCenterX -= ScreenWidth
	}
}

func (u *unit) moveLeft(dist float32) {
	u.headCenterX -= dist

	// teleport if head center is offscreen.
	if u.headCenterX < 0 {
		u.headCenterX += ScreenWidth
	}
}

func (u *unit) markHeadCenter(dst *ebiten.Image) {
	headCX := float64(u.headCenterX)
	headCY := float64(u.headCenterY)
	ebitenutil.DrawLine(dst, headCX-3, headCY, headCX+3, headCY, colorFood)
	ebitenutil.DrawLine(dst, headCX, headCY-3, headCX, headCY+3, colorFood)
}

// Implement collidable interface
// ------------------------------
func (u *unit) collEnabled() bool {
	return true
}

func (u *unit) Rects() []rectF32 {
	return u.rects
}
