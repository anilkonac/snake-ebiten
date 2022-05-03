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

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	foodLength     = 16
	halfFoodLength = foodLength / 2.0
)

type food struct {
	isActive         bool
	centerX, centerY float32
	rects            []rectF32
}

func newFood(centerX, centerY float32) *food {
	newFood := &food{
		centerX: centerX,
		centerY: centerY,
		rects:   make([]rectF32, 0, 4),
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

// Implement collidable interface
// ------------------------------
func (f food) collEnabled() bool {
	return true
}

func (f food) collisionRects() []rectF32 {
	return f.rects
}

// Implement drawable interface
// ------------------------------
func (f food) drawEnabled() bool {
	return f.isActive
}

func (f food) drawableRects() []rectF32 {
	return f.rects
}

func (f food) Color() color.Color {
	return colorFood
}

func (f food) drawingSize() *[2]float32 {
	return &[2]float32{foodLength, foodLength}
}

func (f food) drawDebugInfo(dst *ebiten.Image) {
	cX := float64(f.centerX)
	cY := float64(f.centerY)
	markPoint(dst, cX, cY, 4, colorSnake1)
	for iRect := range f.rects {
		rect := f.rects[iRect]
		rect.drawOuterRect(dst, colorSnake1)
	}
}
