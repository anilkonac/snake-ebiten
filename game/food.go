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
	rects    []rectF64
}

func newFood(centerX, centerY float64) *food {
	newFood := &food{
		rects: make([]rectF64, 0, 4),
	}

	// Create a rectangle to use in drawing and eating logic.
	pureRect := rectF64{
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
	return newFood(float64(rand.Intn(ScreenWidth)), float64(rand.Intn(ScreenHeight)))
}

// Implement slicer interface
// --------------------------
func (f food) slice() []rectF64 {
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

func (f food) totalDimension() (width, height float64) {
	return foodLength, foodLength
}
