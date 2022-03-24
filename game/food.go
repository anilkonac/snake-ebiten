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
	"math/rand"
)

const (
	foodLength     = snakeWidth / 2.0
	halfFoodLength = halfSnakeWidth
)

type food struct {
	isActive bool
	rects    []rectF32
}

func newFood(centerX, centerY float32) *food {
	newFood := &food{
		rects: make([]rectF32, 0, 4),
	}

	// Create a rectangle to use in drawing and eating logic.
	pureRect := rectF32{
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
	return newFood(float32(rand.Intn(ScreenWidth)), float32(rand.Intn(ScreenHeight)))
}

// Implement slicer interface
// --------------------------
func (f food) slice() []rectF32 {
	return f.rects
}

// Implement collidable interface
// ------------------------------
func (f food) collEnabled() bool {
	return true
}

// Implement drawable interface
// ------------------------------
func (f food) drawEnabled() bool {
	return f.isActive
}

func (f food) Color() color.Color {
	return colorFood
}
