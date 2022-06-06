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
	"image/color"

	"github.com/anilkonac/snake-ebiten/game/param"
	"github.com/anilkonac/snake-ebiten/game/shader"
	t "github.com/anilkonac/snake-ebiten/game/tool"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	param.ShaderRound = t.NewShader(shader.Round)
}

type collidable interface {
	CollEnabled() bool
	CollisionRects() []t.RectF32
}

type drawable interface {
	DrawEnabled() bool
	DrawableRects() []t.RectF32
	Color() *color.RGBA
	DrawOptions() *ebiten.DrawTrianglesShaderOptions
	Shader() *ebiten.Shader
	DrawDebugInfo(dst *ebiten.Image)
}

func Draw(dst *ebiten.Image, src drawable) {
	if !src.DrawEnabled() {
		return
	}

	vertices, indices := triangles(src)
	dst.DrawTrianglesShader(vertices, indices, src.Shader(), src.DrawOptions())

	if param.DebugUnits {
		src.DrawDebugInfo(dst)
	}
}

func triangles(src drawable) (vertices []ebiten.Vertex, indices []uint16) {
	vertices = make([]ebiten.Vertex, 0, 16)
	indices = make([]uint16, 0, 24)
	var offset uint16

	rects := src.DrawableRects()
	for iRect := range rects {
		rect := &rects[iRect]

		verticesRect := rect.Vertices(src.Color())
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

func Collides(a, b collidable, tolerance float32) bool {
	if !a.CollEnabled() || !b.CollEnabled() {
		return false
	}

	rectsA := a.CollisionRects()
	rectsB := b.CollisionRects()

	for iRectA := range rectsA {
		rectA := &rectsA[iRectA]

		for iRectB := range rectsB {
			rectB := &rectsB[iRectB]

			if !t.Intersects(rectA, rectB, tolerance) {
				continue
			}

			return true
		}
	}

	return false
}
