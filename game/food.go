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
	t "github.com/anilkonac/snake-ebiten/game/tools"
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
	center   t.Vec32
	rects    []t.RectF32
}

func newFood(center t.Vec32) *food {
	newFood := &food{
		center: center,
		rects:  make([]t.RectF32, 0, 4),
	}

	// Create a rectangle to use in drawing and eating logic.
	pureRect := t.RectF32{
		Pos: t.Vec32{
			X: center.X - params.RadiusFood,
			Y: center.Y - params.RadiusFood,
		},
		Size: t.Vec32{X: params.FoodLength, Y: params.FoodLength},
	}
	// Split this rectangle if it is on a screen edge.
	pureRect.Split(&newFood.rects)

	return newFood
}

func newFoodRandLoc() *food {
	return newFood(t.VecI{X: rand.Intn(params.ScreenWidth), Y: rand.Intn(params.ScreenHeight)}.To32())
}

// Implement collidable interface
// ------------------------------
func (f food) collEnabled() bool {
	return true
}

func (f food) collisionRects() []t.RectF32 {
	return f.rects
}

// Implement drawable interface
// ----------------------------
func (f food) drawEnabled() bool {
	return f.isActive
}

func (f food) drawableRects() []t.RectF32 {
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
	t.MarkPoint(dst, f.center.To64(), 4, params.ColorSnake1)
	for iRect := range f.rects {
		rect := f.rects[iRect]
		rect.DrawOuterRect(dst, params.ColorSnake1)
	}
}
