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

	"github.com/anilkonac/snake-ebiten/game/params"
	"github.com/hajimehoshi/ebiten/v2"
)

var drawOptionsFood = ebiten.DrawTrianglesShaderOptions{
	Uniforms: map[string]interface{}{
		"Radius":     float32(params.RadiusFood),
		"IsVertical": float32(1.0),
		"Size":       []float32{params.FoodLength, params.FoodLength},
	},
}

type food struct {
	isActive bool
	center   vec32
	rects    []rectF32
}

func newFood(center vec32) *food {
	newFood := &food{
		center: center,
		rects:  make([]rectF32, 0, 4),
	}

	// Create a rectangle to use in drawing and eating logic.
	pureRect := rectF32{
		pos: vec32{
			x: center.x - params.RadiusFood,
			y: center.y - params.RadiusFood,
		},
		size: vec32{params.FoodLength, params.FoodLength},
	}
	// Split this rectangle if it is on a screen edge.
	pureRect.split(&newFood.rects)

	return newFood
}

func newFoodRandLoc() *food {
	return newFood(vecI{rand.Intn(params.ScreenWidth), rand.Intn(params.ScreenHeight)}.to32())
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
// ----------------------------
func (f food) drawEnabled() bool {
	return f.isActive
}

func (f food) drawableRects() []rectF32 {
	return f.rects
}

func (f food) Color() *color.RGBA {
	return &params.ColorFood
}

func (f food) drawOptions() *ebiten.DrawTrianglesShaderOptions {
	return &drawOptionsFood
}

func (f food) shader() *ebiten.Shader {
	return params.ShaderRound
}

func (f food) drawDebugInfo(dst *ebiten.Image) {
	markPoint(dst, f.center.to64(), 4, params.ColorSnake1)
	for iRect := range f.rects {
		rect := f.rects[iRect]
		rect.drawOuterRect(dst, params.ColorSnake1)
	}
}
