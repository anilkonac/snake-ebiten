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

	"github.com/hajimehoshi/ebiten/v2"
)

type slicer interface {
	slice() []rect
}

type collidable interface {
	slicer
	collEnabled() bool
}

type drawable interface {
	slicer
	drawEnabled() bool
	Color() color.Color
	dimension() *[2]float32
	drawDebugInfo(dst *ebiten.Image)
}

func draw(dst *ebiten.Image, src drawable) {
	if !src.drawEnabled() {
		return
	}

	var radius float32
	var isVertical float32
	switch v := src.(type) {
	case *unit:
		radius = halfSnakeWidth
		if v.direction.isVertical() {
			isVertical = 1.0
		}
	case *food:
		radius = halfFoodLength
	}

	vertices, indices := triangles(src)
	op := &ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]interface{}{
			"Radius":     radius,
			"IsVertical": isVertical,
			"Dimension":  (*src.dimension())[:],
		},
	}
	dst.DrawTrianglesShader(vertices, indices, shaderMap[curShader], op)

	if debugUnits {
		src.drawDebugInfo(dst)
	}
}

func triangles(src drawable) (vertices []ebiten.Vertex, indices []uint16) {
	vertices = make([]ebiten.Vertex, 0, 16)
	indices = make([]uint16, 0, 24)
	var offset uint16

	rects := src.slice()
	for iRect := range rects {
		rect := &rects[iRect]

		verticesRect := rect.vertices(src.Color())
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

func collides(a, b collidable, tolerance int16) bool {
	if !a.collEnabled() || !b.collEnabled() {
		return false
	}

	rectsA := a.slice()
	rectsB := b.slice()

	for iRectA := range rectsA {
		rectA := &rectsA[iRectA]
		for iRectB := range rectsB {
			rectB := &rectsB[iRectB]
			if !intersects(rectA, rectB, tolerance) {
				continue
			}

			return true
		}
	}

	return false
}
