// snake-ebiten
// Copyright (C) 2022 Anıl Konaç

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

//go:build ignore

package main

var Color vec4
var ShadedCorners [4]float
var RectSize vec2
var RectPosInUnit vec2
var TotalSize vec2
var IsVertical float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	normColor := Color / 0xffff

	posInUnit := vec2(RectSize.x*texCoord.x, RectSize.y*texCoord.y) + RectPosInUnit
	if IsVertical > 0 {
		radius := TotalSize.x / 2

		// Top Left Corner
		if (ShadedCorners[0] > 0) && (posInUnit.x < radius) &&
			(posInUnit.y < radius) && (posInUnit.y < growUp(posInUnit.x, radius, TotalSize.y)) {
			normColor.a = 0
		}
		// Bottom Left Corner
		if (ShadedCorners[1] > 0) && (posInUnit.x < radius) &&
			(posInUnit.y > (TotalSize.y - radius)) && (posInUnit.y > growDown(posInUnit.x, radius, TotalSize.y)) {
			normColor.a = 0
		}
		// Bottom Right Corner
		if (ShadedCorners[2] > 0) && (posInUnit.x > radius) &&
			(posInUnit.y > (TotalSize.y - radius)) && (posInUnit.y > growDown(posInUnit.x, radius, TotalSize.y)) {
			normColor.a = 0
		}
		// Top Right Corner
		if (ShadedCorners[3] > 0) && (posInUnit.x > radius) &&
			(posInUnit.y < radius) && (posInUnit.y < growUp(posInUnit.x, radius, TotalSize.y)) {
			normColor.a = 0
		}

	} else {
		// if TotalSize.x >= TotalSize.y {
		radius := TotalSize.y / 2
		roundCenter1 := vec2(radius)
		roundCenter2 := vec2(TotalSize.x-radius, radius)

		distToCenter1 := distance(posInUnit, roundCenter1)
		distToCenter2 := distance(posInUnit, roundCenter2)

		// Top Left Corner
		if (ShadedCorners[0] > 0) && (posInUnit.x < radius) &&
			(posInUnit.y < radius) && (distToCenter1 > radius) {
			normColor.a = 0
		}
		// Bottom Left Corner
		if (ShadedCorners[1] > 0) && (posInUnit.x < radius) &&
			(posInUnit.y > radius) && (distToCenter1 > radius) {
			normColor.a = 0
		}
		// Bottom Right Corner
		if (ShadedCorners[2] > 0) && (posInUnit.x > TotalSize.x-radius) &&
			(posInUnit.y > radius) && (distToCenter2 > radius) {
			normColor.a = 0
		}
		// Top Right Corner
		if (ShadedCorners[3] > 0) && (posInUnit.x > TotalSize.x-radius) &&
			(posInUnit.y < radius) && (distToCenter2 > radius) {
			normColor.a = 0
		}
	}

	normColor.rgb *= normColor.a
	return normColor
}

func growUp(x, radius, height float) float {
	// Linear interpolation between squared function and semicircle function
	heightMultip := clamp(TotalSize.y/radius, 0, 1)
	y0 := heightMultip * (pow(x-radius, 2) / radius)
	x0 := 0.0
	x1 := 1.0
	y1 := radius - sqrt(radius*radius-(x-radius)*(x-radius))

	return y0 + (heightMultip-x0)*(y1-y0)/(x1-x0)
}

func growDown(x, radius, height float) float {
	// Linear interpolation between squared function and semicircle function
	heightMultip := clamp(TotalSize.y/radius, 0, 1)
	y0 := heightMultip * (0.5 - pow(x-radius, 2)/radius)
	x0 := 0.0
	x1 := 1.0
	y1 := height - heightMultip*radius + sqrt(pow(radius, 2)-pow(x-radius, 2))

	return y0 + (heightMultip-x0)*(y1-y0)/(x1-x0)
}
