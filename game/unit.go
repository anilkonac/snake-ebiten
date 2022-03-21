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
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type unit struct {
	headCenterX float64
	headCenterY float64
	length      float64
	direction   directionT
	rects       []rectF64 // rectangles that are used for both collision checking and drawing
	color       *color.RGBA
	next        *unit
	prev        *unit
}

func newUnit(headCenterX, headCenterY, length float64, direction directionT, color *color.RGBA) *unit {
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
	var pureRect rectF64
	length64 := float64(u.length)
	switch u.direction {
	case directionRight:
		pureRect = rectF64{
			x:      u.headCenterX - length64 + halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  length64,
			height: snakeWidth,
		}
	case directionLeft:
		pureRect = rectF64{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  length64,
			height: snakeWidth,
		}
	case directionUp:
		pureRect = rectF64{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  snakeWidth,
			height: length64,
		}
	case directionDown:
		pureRect = rectF64{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - length64 + halfSnakeWidth,
			width:  snakeWidth,
			height: length64,
		}
	default:
		panic("Wrong unit direction!!")
	}

	u.rects = make([]rectF64, 0, 4) // Remove old rectangles
	pureRect.split(&u.rects)        // Create split rectangles on screen edges.
}

func (u *unit) moveUp(dist float64) {
	u.headCenterY -= dist

	// teleport if head center is offscreen.
	if u.headCenterY < 0 {
		u.headCenterY += ScreenHeight
	}
}

func (u *unit) moveDown(dist float64) {
	u.headCenterY += dist

	// teleport if head center is offscreen.
	if u.headCenterY > ScreenHeight {
		u.headCenterY -= ScreenHeight
	}
}

func (u *unit) moveRight(dist float64) {
	u.headCenterX += dist

	// teleport if head center is offscreen.
	if u.headCenterX > ScreenWidth {
		u.headCenterX -= ScreenWidth
	}
}

func (u *unit) moveLeft(dist float64) {
	u.headCenterX -= dist

	// teleport if head center is offscreen.
	if u.headCenterX < 0 {
		u.headCenterX += ScreenWidth
	}
}

func (u *unit) markHeadCenter(dst *ebiten.Image) {
	headCX := u.headCenterX
	headCY := u.headCenterY
	ebitenutil.DrawLine(dst, headCX-3, headCY, headCX+3, headCY, colorFood)
	ebitenutil.DrawLine(dst, headCX, headCY-3, headCX, headCY+3, colorFood)
}

// Compute rounding of corners for body units
func (u *unit) roundCornersBody(roundCorners *[4]float32, headUnit *unit) {
	headGrowth := float32(math.Min(headUnit.length/halfSnakeWidth, 1.0))
	roundMult := 1 - headGrowth
	switch {
	case u.prev.direction == directionUp && u.direction == directionRight:
		roundCorners[2] = 1
		if u.prev == headUnit {
			roundCorners[3] = roundMult
		}
	case u.prev.direction == directionLeft && u.direction == directionUp:
		roundCorners[3] = 1
		if u.prev == headUnit {
			roundCorners[0] = roundMult
		}
	case u.prev.direction == directionDown && u.direction == directionLeft:
		roundCorners[0] = 1
		if u.prev == headUnit {
			roundCorners[1] = roundMult
		}
	case u.prev.direction == directionRight && u.direction == directionDown:
		roundCorners[1] = 1
		if u.prev == headUnit {
			roundCorners[2] = roundMult
		}
	case u.prev.direction == directionDown && u.direction == directionRight:
		roundCorners[3] = 1
		if u.prev == headUnit {
			roundCorners[2] = roundMult
		}
	case u.prev.direction == directionLeft && u.direction == directionDown:
		roundCorners[2] = 1
		if u.prev == headUnit {
			roundCorners[1] = roundMult
		}
	case u.prev.direction == directionUp && u.direction == directionLeft:
		roundCorners[1] = 1
		if u.prev == headUnit {
			roundCorners[0] = roundMult
		}
	case u.prev.direction == directionRight && u.direction == directionUp:
		roundCorners[0] = 1
		if u.prev == headUnit {
			roundCorners[3] = roundMult
		}
	}
	return
}

// Compute rounding of corners of the unit next to the tail
func (u *unit) roundCornersPreTail(roundCorners *[4]float32, tailUnit *unit) {
	// if next unit is tail and its length is less than half snake width
	if (u.next == tailUnit) && (tailUnit.length < snakeWidth) {
		tailShrink := float32(-1.0 + tailUnit.length/halfSnakeWidth)
		switch {
		case (u.direction == directionUp) && (tailUnit.direction == directionRight):
			roundCorners[1] = tailShrink
			roundCorners[2] = 1.0
		case (u.direction == directionUp) && (tailUnit.direction == directionLeft):
			roundCorners[1] = 1.0
			roundCorners[2] = tailShrink
		case (u.direction == directionDown) && (tailUnit.direction == directionLeft):
			roundCorners[0] = 1.0
			roundCorners[3] = tailShrink
		case (u.direction == directionDown) && (tailUnit.direction == directionRight):
			roundCorners[0] = tailShrink
			roundCorners[3] = 1.0
		case (u.direction == directionRight) && (tailUnit.direction == directionUp):
			roundCorners[0] = 1.0
			roundCorners[1] = tailShrink
		case (u.direction == directionRight) && (tailUnit.direction == directionDown):
			roundCorners[0] = tailShrink
			roundCorners[1] = 1.0
		case (u.direction == directionLeft) && (tailUnit.direction == directionUp):
			roundCorners[2] = tailShrink
			roundCorners[3] = 1.0
		case (u.direction == directionLeft) && (tailUnit.direction == directionDown):
			roundCorners[2] = 1.0
			roundCorners[3] = tailShrink
		}
	}
}

// Implement slicer interface
// --------------------------
func (u *unit) slice() []rectF64 {
	return u.rects
}

// Implement collidable interface
// ------------------------------
func (u *unit) collEnabled() bool {
	return true
}

// Implement drawable interface
// ------------------------------
func (u *unit) drawEnabled() bool {
	return true
}

func (u *unit) Color() color.Color {
	return u.color
}

func (u *unit) totalDimension() *[2]float64 {
	if u.direction.isVertical() {
		return &[2]float64{snakeWidth, u.length}
	}
	return &[2]float64{u.length, snakeWidth}
}

// func (u *unit) shader() shaderT {
// 	return shaderRound
// }
