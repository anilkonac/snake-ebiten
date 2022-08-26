/*
Copyright (C) 2022 Anıl Konaç

This file is part of snake-ebiten.

snake-ebiten is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

snake-ebiten is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with snake-ebiten. If not, see <https://www.gnu.org/licenses/>.
*/

package object

import (
	"math/rand"

	c "github.com/anilkonac/snake-ebiten/game/core"
	"github.com/anilkonac/snake-ebiten/game/param"
	"github.com/anilkonac/snake-ebiten/game/shader"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	imageFood    = ebiten.NewImage(param.FoodLength, param.FoodLength)
	foodDrawOpts ebiten.DrawTrianglesOptions
)

func init() {
	imageFood.DrawRectShader(param.FoodLength, param.FoodLength, &shader.Circle, &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]interface{}{
			"Radius": float32(param.RadiusFood),
		},
	})
}

type Food struct {
	c.TeleCompTriang
	IsActive bool
	Center   c.Vec32
}

func newFood(center c.Vec32) *Food {
	newFood := &Food{
		Center: center,
	}
	newFood.SetColor(&param.ColorFood)

	// Create a rectangle to use in drawing and eating logic.
	pureRect := c.RectF32{
		Pos: c.Vec32{
			X: center.X - param.RadiusFood,
			Y: center.Y - param.RadiusFood,
		},
		Size: c.Vec32{X: param.FoodLength, Y: param.FoodLength},
	}
	// Split this rectangle if it is on a screen edge.
	newFood.Update(&pureRect)

	return newFood
}

func NewFoodRandLoc() *Food {
	return newFood(c.VecI{X: rand.Intn(param.ScreenWidth), Y: rand.Intn(param.ScreenHeight)}.To32())
}

func (f Food) Draw(dst *ebiten.Image) {
	vertices, indices := f.TeleCompTriang.Triangles()
	dst.DrawTriangles(vertices, indices, imageFood, &foodDrawOpts)
}

// Implement collidable interface
// ------------------------------
func (f Food) CollEnabled() bool {
	return true
}

func (f Food) CollisionRects() []c.RectF32 {
	return f.Rects[:]
}
