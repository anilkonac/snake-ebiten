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
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	foodLength     = 16
	halfFoodLength = foodLength / 2.0
)

type food struct {
	isActive         bool
	centerX, centerY int16 // for debugging purposes
	rects            []rect
}

func newFood(centerX, centerY int16) *food {
	newFood := &food{
		centerX: centerX,
		centerY: centerY,
		rects:   make([]rect, 0, 4),
	}

	// Create a rectangle to use in drawing and eating logic.
	pureRect := rect{
		x:      centerX - halfFoodLength,
		y:      centerY - halfFoodLength,
		width:  foodLength,
		height: foodLength,
	}
	// Split this rectangle if it is on a screen edge.
	pureRect.split(&newFood.rects)

	return newFood
}

func newFoodRandLoc() *food {
	return newFood(int16(rand.Intn(ScreenWidth)), int16(rand.Intn(ScreenHeight)))
}

// Implement collidable interface
// ------------------------------
func (f food) collEnabled() bool {
	return true
}

func (f food) Rects() []rect {
	return f.rects
}

// Implement drawable interface
// ------------------------------
func (f food) drawEnabled() bool {
	return f.isActive
}

func (f food) triangles() (vertices []ebiten.Vertex, indices []uint16) {
	vertices = make([]ebiten.Vertex, 0)
	indices = make([]uint16, 0)
	var offset uint16

	for iRect := range f.rects {
		rect := &f.rects[iRect]

		verticesRect := rect.vertices(colorFood)
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

func (f food) dimension() *[2]float32 {
	return &[2]float32{foodLength, foodLength}
}

func (f food) drawDebugInfo(dst *ebiten.Image) {
	cX := float64(f.centerX)
	cY := float64(f.centerY)
	markPoint(dst, cX, cY, colorSnake1)
	for iRect := range f.rects {
		rect := f.rects[iRect]
		rect.drawOuterRect(dst, colorSnake1)
	}
}
