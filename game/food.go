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

func (f food) draw(dst *ebiten.Image) {
	if !f.isActive {
		return
	}

	vertices, indices := f.triangles()
	op := &ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]interface{}{
			"Radius":     float32(halfFoodLength),
			"IsVertical": float32(1.0),
			"Dimension":  []float32{foodLength, foodLength},
		},
	}
	dst.DrawTrianglesShader(vertices, indices, shaderMap[shaderRound], op)
}

// Implement collidable interface
// ------------------------------
func (f food) collEnabled() bool {
	return true
}

func (f food) Rects() []rectF32 {
	return f.rects
}

// // Implement drawable interface
// // ------------------------------
// func (f food) drawEnabled() bool {
// 	return f.isActive
// }

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

// func (f food) drawDebugInfo(dst *ebiten.Image) {
// 	for iRect := range f.rects {
// 		rect := &f.rects[iRect]

// 		ebitenutil.DebugPrintAt(dst, fmt.Sprintf("%3.3f, %3.3f", rect.x, rect.y), int(rect.x)-90, int(rect.y)-15)
// 		bottomX := rect.x + rect.width
// 		bottomY := rect.y + rect.height
// 		ebitenutil.DebugPrintAt(dst, fmt.Sprintf("%3.3f, %3.3f", bottomX, bottomY), int(bottomX), int(bottomY))
// 	}
// }
