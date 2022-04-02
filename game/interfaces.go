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

type collidable interface {
	collEnabled() bool
	Rects() []rectF32
}

// type drawable interface {
// 	drawEnabled() bool
// 	triangles() (vertices []ebiten.Vertex, indices []uint16)
// 	drawDebugInfo(dst *ebiten.Image)
// }

// func draw(dst *ebiten.Image, src drawable) {
// 	if !src.drawEnabled() {
// 		return
// 	}

// 	vertices, indices := src.triangles()
// 	dst.DrawTrianglesShader(vertices, indices, shaderMap[shaderBasic], new(ebiten.DrawTrianglesShaderOptions))

// 	if debugUnits {
// 		src.drawDebugInfo(dst)
// 	}
// }

func collides(a, b collidable, tolerance float32) bool {
	if !a.collEnabled() || !b.collEnabled() {
		return false
	}

	rectsA := a.Rects()
	rectsB := b.Rects()

	for iRectA := range rectsA {
		rectA := &rectsA[iRectA]
		for iRectB := range b.Rects() {
			rectB := &rectsB[iRectB]
			if !intersects(rectA, rectB, tolerance) {
				continue
			}

			return true
		}
	}

	return false
}
