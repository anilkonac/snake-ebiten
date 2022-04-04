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
	headCenterX float64
	headCenterY float64
	length      float64
	direction   directionT
	rects       []rect // rectangles that are used for both collision checking and drawing
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
	var pureRect rect
	iLength := int16(u.length)
	iCX := int16(u.headCenterX)
	iCY := int16(u.headCenterY)
	switch u.direction {
	case directionRight:
		pureRect = rect{
			x:      iCX - iLength + halfSnakeWidth,
			y:      iCY - halfSnakeWidth,
			width:  iLength,
			height: snakeWidth,
		}
	case directionLeft:
		pureRect = rect{
			x:      iCX - halfSnakeWidth,
			y:      iCY - halfSnakeWidth,
			width:  iLength,
			height: snakeWidth,
		}
	case directionUp:
		pureRect = rect{
			x:      iCX - halfSnakeWidth,
			y:      iCY - halfSnakeWidth,
			width:  snakeWidth,
			height: iLength,
		}
	case directionDown:
		pureRect = rect{
			x:      iCX - halfSnakeWidth,
			y:      iCY - iLength + halfSnakeWidth,
			width:  snakeWidth,
			height: iLength,
		}
	default:
		panic("Wrong unit direction!!")
	}

	u.rects = make([]rect, 0, 4) // Remove old rectangles
	pureRect.split(&u.rects)     // Create split rectangles on screen edges.
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

	switch u.direction {
	case directionUp:
		headCY = float64(u.headCenterY+u.length) - snakeWidth
	case directionDown:
		headCY = float64(u.headCenterY-u.length) + snakeWidth
	case directionRight:
		headCX = float64(u.headCenterX-u.length) + snakeWidth
	case directionLeft:
		headCX = float64(u.headCenterX+u.length) - snakeWidth
	}
	// mark head center at the other side
	markPoint(dst, headCX, headCY, colorFood)
}

// Implement collidable interface
// ------------------------------
func (u *unit) collEnabled() bool {
	return true
}

func (u *unit) Rects() []rect {
	return u.rects
}

// Implement drawable interface
// ----------------------------
func (u *unit) drawEnabled() bool {
	return true
}

func (u *unit) triangles() (vertices []ebiten.Vertex, indices []uint16) {
	vertices = make([]ebiten.Vertex, 0, 16)
	indices = make([]uint16, 0, 24)
	var offset uint16

	for iRect := range u.rects {
		rect := &u.rects[iRect]

		verticesRect := rect.vertices(u.color)
		indicesRect := []uint16{
			offset + 1, offset, offset + 2,
			offset + 2, offset + 3, offset + 1,
		}

		vertices = append(vertices, verticesRect...)
		indices = append(indices, indicesRect...)

		offset += 4
	}

	return
}

func (u *unit) dimension() *[2]float32 {
	flooredLength := float32(math.Floor(u.length))
	if u.direction.isVertical() {
		return &[2]float32{snakeWidth, flooredLength}
	}
	return &[2]float32{flooredLength, snakeWidth}
}

func (u *unit) drawDebugInfo(dst *ebiten.Image) {
	u.markHeadCenters(dst)
	for iRect := range u.rects {
		rect := u.rects[iRect]
		rect.drawOuterRect(dst, colorFood)
	}
}
