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
)

type unit struct {
	headCenterX    float64
	headCenterY    float64
	length         float64
	direction      directionT
	rectsCollision []rectF32
	rectsDrawable  []rectF32
	color          *color.RGBA
	next           *unit
	prev           *unit
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
	// Create rectangles for drawing and collision. They are going to split.
	var rectDraw, rectColl *rectF32
	length32 := float32(math.Floor(u.length))
	cX32 := float32(math.Floor(u.headCenterX))
	cY32 := float32(math.Floor(u.headCenterY))
	switch u.direction {
	case directionRight:
		rectColl = newRect(cX32-length32+halfSnakeWidth, cY32-halfSnakeWidth, length32, snakeWidth)
		if u.next != nil {
			rectDraw = newRect(rectColl.x-snakeWidth, rectColl.y, rectColl.width+snakeWidth, rectColl.height)
		}
	case directionLeft:
		rectColl = newRect(cX32-halfSnakeWidth, cY32-halfSnakeWidth, length32, snakeWidth)
		if u.next != nil {
			rectDraw = newRect(rectColl.x, rectColl.y, rectColl.width+snakeWidth, rectColl.height)
		}
	case directionUp:
		rectColl = newRect(cX32-halfSnakeWidth, cY32-halfSnakeWidth, snakeWidth, length32)
		if u.next != nil {
			rectDraw = newRect(rectColl.x, rectColl.y, rectColl.width, rectColl.height+snakeWidth)
		}
	case directionDown:
		rectColl = newRect(cX32-halfSnakeWidth, cY32-length32+halfSnakeWidth, snakeWidth, length32)
		if u.next != nil {
			rectDraw = newRect(rectColl.x, rectColl.y-snakeWidth, rectColl.width, rectColl.height+snakeWidth)
		}
	default:
		panic("Wrong unit direction!!")
	}

	// Remove old rectangles
	u.rectsDrawable = make([]rectF32, 0, 4)
	u.rectsCollision = make([]rectF32, 0, 4)

	// Create split rectangles on screen edges.
	rectColl.split(&u.rectsCollision)
	if u.next == nil {
		u.rectsDrawable = u.rectsCollision
		return
	}
	rectDraw.split(&u.rectsDrawable)
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

func (u *unit) markHeadCenters(dst *ebiten.Image) {
	headCX := float64(u.headCenterX)
	headCY := float64(u.headCenterY)
	markPoint(dst, headCX, headCY, colorFood)

	var offset float64 = 0
	if u.next == nil {
		offset = snakeWidth
	}
	switch u.direction {
	case directionUp:
		headCY = u.headCenterY + u.length - offset
	case directionDown:
		headCY = float64(u.headCenterY-u.length) + offset
	case directionRight:
		headCX = float64(u.headCenterX-u.length) + offset
	case directionLeft:
		headCX = float64(u.headCenterX+u.length) - offset
	}
	// mark head center at the other side
	markPoint(dst, headCX, headCY, colorFood)
}

// Implement collidable interface
// ------------------------------
func (u *unit) collEnabled() bool {
	return true
}

func (u *unit) collisionRects() []rectF32 {
	return u.rectsCollision
}

// Implement drawable interface
// ----------------------------
func (u *unit) drawEnabled() bool {
	return true
}

func (u *unit) drawableRects() []rectF32 {
	return u.rectsDrawable
}

func (u *unit) Color() color.Color {
	return u.color
}

func (u *unit) drawDimension() *[2]float32 {
	flooredLength := float32(math.Floor(u.length))
	if u.next != nil {
		flooredLength += snakeWidth
	}
	if u.direction.isVertical() {
		return &[2]float32{snakeWidth, flooredLength}
	}
	return &[2]float32{flooredLength, snakeWidth}
}

func (u *unit) drawDebugInfo(dst *ebiten.Image) {
	u.markHeadCenters(dst)
	for iRect := range u.rectsDrawable {
		rect := u.rectsDrawable[iRect]
		rect.drawOuterRect(dst, colorFood)
	}
}
