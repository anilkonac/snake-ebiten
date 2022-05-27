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

package object

import (
	"image/color"
	"math/rand"

	"github.com/anilkonac/snake-ebiten/game/param"
	t "github.com/anilkonac/snake-ebiten/game/tool"
	"github.com/hajimehoshi/ebiten/v2"
)

var drawOptionsFood = ebiten.DrawTrianglesShaderOptions{
	Uniforms: map[string]interface{}{
		"Radius":     float32(param.RadiusFood),
		"IsVertical": float32(1.0),
		"Size":       []float32{param.FoodLength, param.FoodLength},
	},
}

type Food struct {
	IsActive bool
	Center   t.Vec32
	rects    []t.RectF32
}

func newFood(center t.Vec32) *Food {
	newFood := &Food{
		Center: center,
		rects:  make([]t.RectF32, 0, 4),
	}

	// Create a rectangle to use in drawing and eating logic.
	pureRect := t.RectF32{
		Pos: t.Vec32{
			X: center.X - param.RadiusFood,
			Y: center.Y - param.RadiusFood,
		},
		Size: t.Vec32{X: param.FoodLength, Y: param.FoodLength},
	}
	// Split this rectangle if it is on a screen edge.
	pureRect.Split(&newFood.rects)

	return newFood
}

func NewFoodRandLoc() *Food {
	return newFood(t.VecI{X: rand.Intn(param.ScreenWidth), Y: rand.Intn(param.ScreenHeight)}.To32())
}

// Implement collidable interface
// ------------------------------
func (f Food) CollEnabled() bool {
	return true
}

func (f Food) CollisionRects() []t.RectF32 {
	return f.rects
}

// Implement drawable interface
// ----------------------------
func (f Food) DrawEnabled() bool {
	return f.IsActive
}

func (f Food) DrawableRects() []t.RectF32 {
	return f.rects
}

func (f Food) Color() *color.RGBA {
	return &param.ColorFood
}

func (f Food) DrawOptions() *ebiten.DrawTrianglesShaderOptions {
	return &drawOptionsFood
}

func (f Food) Shader() *ebiten.Shader {
	return param.ShaderRound
}

func (f Food) DrawDebugInfo(dst *ebiten.Image) {
	t.MarkPoint(dst, f.Center.To64(), 4, param.ColorSnake1)
	for iRect := range f.rects {
		rect := f.rects[iRect]
		rect.DrawOuterRect(dst, param.ColorSnake1)
	}
}
