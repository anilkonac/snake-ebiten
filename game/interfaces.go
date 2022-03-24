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
	slice() []rectF32
}

type collidable interface {
	slicer
	collEnabled() bool
}

type drawable interface {
	slicer
	drawEnabled() bool
	Color() color.Color
}

func draw(dst *ebiten.Image, src drawable) {
	if !src.drawEnabled() {
		return
	}

	for _, rect := range src.slice() {
		rect.draw(dst, src.Color())
	}
}

func collides(a, b collidable, tolerance float32) bool {
	if !a.collEnabled() || !b.collEnabled() {
		return false
	}

	for _, rectA := range a.slice() {
		for _, rectB := range b.slice() {
			if !intersects(rectA, rectB, tolerance) {
				continue
			}

			return true
		}
	}

	return false
}
