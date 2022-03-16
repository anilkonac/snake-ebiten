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
			(posInUnit.y < radius) && (posInUnit.y < growUp(posInUnit.x, radius)) {
			normColor.a = 0
		}
		// Bottom Left Corner
		if (ShadedCorners[1] > 0) && (posInUnit.x < radius) &&
			(posInUnit.y > (TotalSize.y - radius)) && (posInUnit.y > growDown(posInUnit.x, radius)) {
			normColor.a = 0
		}
		// Bottom Right Corner
		if (ShadedCorners[2] > 0) && (posInUnit.x > radius) &&
			(posInUnit.y > (TotalSize.y - radius)) && (posInUnit.y > growDown(posInUnit.x, radius)) {
			normColor.a = 0
		}
		// Top Right Corner
		if (ShadedCorners[3] > 0) && (posInUnit.x > radius) &&
			(posInUnit.y < radius) && (posInUnit.y < growUp(posInUnit.x, radius)) {
			normColor.a = 0
		}

	} else {
		// if TotalSize.x >= TotalSize.y {
		radius := TotalSize.y / 2

		// Top Left Corner
		if (ShadedCorners[0] > 0) && (posInUnit.x < radius) &&
			(posInUnit.y < radius) && (posInUnit.x < growLeft(posInUnit.y, radius)) {
			normColor.a = 0
		}
		// Bottom Left Corner
		if (ShadedCorners[1] > 0) && (posInUnit.x < radius) &&
			(posInUnit.y > radius) && (posInUnit.x < growLeft(posInUnit.y, radius)) {
			normColor.a = 0
		}
		// Bottom Right Corner
		if (ShadedCorners[2] > 0) && (posInUnit.x > TotalSize.x-radius) &&
			(posInUnit.y > radius) && (posInUnit.x > growRight(posInUnit.y, radius)) {
			normColor.a = 0
		}
		// Top Right Corner
		if (ShadedCorners[3] > 0) && (posInUnit.x > TotalSize.x-radius) &&
			(posInUnit.y < radius) && (posInUnit.x > growRight(posInUnit.y, radius)) {
			normColor.a = 0
		}
	}

	normColor.rgb *= normColor.a
	return normColor
}

func growUp(x, radius float) float {
	// Linear interpolation between square function and semicircle function
	heightMultip := clamp(TotalSize.y/radius, 0.0, 1.0)
	square := squareGrowthReverse(x, radius, heightMultip)
	semicircle := semicircleReverse(x, radius, heightMultip)

	return heightMultip*semicircle + (1.0-heightMultip)*square
}

func growDown(x, radius float) float {
	// Linear interpolation between square function and semicircle function
	heightMultip := clamp(TotalSize.y/radius, 0.0, 1.0)
	square := squareGrowth(x, radius, heightMultip)
	semicircle := TotalSize.y - radius + semicircle(x, radius, heightMultip)

	return heightMultip*semicircle + (1.0-heightMultip)*square
}

func growLeft(y, radius float) float {
	// Linear interpolation between square function and semicircle function
	widthMultip := clamp(TotalSize.x/radius, 0.0, 1.0)
	square := squareGrowthReverse(y, radius, widthMultip)
	semicircle := semicircleReverse(y, radius, widthMultip)

	return widthMultip*semicircle + (1.0-widthMultip)*square
}

func growRight(y, radius float) float {
	// Linear interpolation between square function and semicircle function
	widthMultip := clamp(TotalSize.x/radius, 0.0, 1.0)
	square := squareGrowth(y, radius, widthMultip)
	semicircle := TotalSize.x - radius + semicircle(y, radius, widthMultip)

	return widthMultip*semicircle + (1.0-widthMultip)*square
}

func squareGrowth(x, radius, multiplier float) float {
	return multiplier * (radius - pow(x-radius, 2.0)/radius)
}

func squareGrowthReverse(x, radius, multiplier float) float {
	return multiplier * pow(x-radius, 2.0) / radius
}

func semicircle(x, radius, multiplier float) float {
	return sqrt(pow(radius, 2.0) - pow(x-radius, 2.0))
}

func semicircleReverse(x, radius, multiplier float) float {
	return radius - semicircle(x, radius, multiplier)
}
