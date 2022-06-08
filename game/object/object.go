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
	"github.com/anilkonac/snake-ebiten/game/core"
	c "github.com/anilkonac/snake-ebiten/game/core"
	"github.com/anilkonac/snake-ebiten/game/param"
	"github.com/anilkonac/snake-ebiten/game/shader"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	param.ShaderRound = c.NewShader(shader.Round)
}

type collidable interface {
	CollEnabled() bool
	CollisionRects() []core.RectF32
}

type drawable interface {
	DrawEnabled() bool
	Triangles() (vertices []ebiten.Vertex, indices []uint16)
	DrawOptions() *ebiten.DrawTrianglesShaderOptions
	Shader() *ebiten.Shader
	DrawDebugInfo(dst *ebiten.Image)
}

func Draw(dst *ebiten.Image, src drawable) {
	if !src.DrawEnabled() {
		return
	}

	vertices, indices := src.Triangles()
	dst.DrawTrianglesShader(vertices, indices, src.Shader(), src.DrawOptions())

	if param.DebugUnits {
		src.DrawDebugInfo(dst)
	}
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

			if !intersects(rectA, rectB, tolerance) {
				continue
			}

			return true
		}
	}

	return false
}

func intersects(rectA, rectB *core.RectF32, tolerance float32) bool {
	aRightX := rectA.Pos.X + rectA.Size.X
	bRightX := rectB.Pos.X + rectB.Size.X
	aBottomY := rectA.Pos.Y + rectA.Size.Y
	bBottomY := rectB.Pos.Y + rectB.Size.Y

	if (rectA.Pos.X-rectB.Pos.X <= tolerance) && (aRightX-rectB.Pos.X <= tolerance) { // rectA is on the left side of rectB
		return false
	}

	if (rectA.Pos.X-bRightX >= -tolerance) && (aRightX-bRightX >= -tolerance) { // rectA is on the right side of rectB
		return false
	}

	if (rectA.Pos.Y-rectB.Pos.Y <= tolerance) && (aBottomY-rectB.Pos.Y <= tolerance) { // rectA is above rectB
		return false
	}

	if (rectA.Pos.Y-bBottomY >= -tolerance) && (aBottomY-bBottomY >= -tolerance) { // rectA is under rectB
		return false
	}

	return true
}
