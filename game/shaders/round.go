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
		radius := TotalSize.x / 2.0

		// Top Left Corner
		if roundMult := ShadedCorners[0]; (roundMult != 0.0) && (posInUnit.x < radius) && (posInUnit.y < radius) { // In point is on top left corner
			multiplier := 1.0
			if (roundMult > 1.0) && (roundMult <= 2) {
				multiplier = 1.0 + (sqrt(2.0)-1.0)*easeOutCirc(roundMult-1.0)
			} else if roundMult < 0 {
				multiplier = sqrt(2) - (sqrt(2.0)-1.0)*easeInSine(-roundMult)
			}
			if posInUnit.y < growUp(posInUnit.x, radius, multiplier) {
				normColor.a = 0.0
			}
		}
		// Bottom Left Corner
		if roundMult := ShadedCorners[1]; (roundMult != 0.0) && (posInUnit.x < radius) && (posInUnit.y > (TotalSize.y - radius)) {
			multiplier := 1.0
			if (roundMult > 1.0) && (roundMult <= 2) {
				multiplier = 1.0 + (sqrt(2.0)-1.0)*easeOutCirc(roundMult-1.0)
			} else if roundMult < 0 {
				multiplier = sqrt(2) - (sqrt(2.0)-1.0)*easeInSine(-roundMult)
			}
			if posInUnit.y > growDown(posInUnit.x, radius, multiplier) {
				normColor.a = 0.0
			}
		}
		// Bottom Right Corner
		if roundMult := ShadedCorners[2]; (roundMult != 0.0) && (posInUnit.x > radius) && (posInUnit.y > (TotalSize.y - radius)) {
			multiplier := 1.0
			if (roundMult > 1.0) && (roundMult <= 2) {
				multiplier = 1.0 + (sqrt(2.0)-1.0)*easeOutCirc(roundMult-1.0)
			} else if roundMult < 0 {
				multiplier = sqrt(2) - (sqrt(2.0)-1.0)*easeInSine(-roundMult)
			}
			if posInUnit.y > growDown(posInUnit.x, radius, multiplier) {
				normColor.a = 0.0
			}
		}
		// Top Right Corner
		if roundMult := ShadedCorners[3]; (roundMult != 0.0) && (posInUnit.x > radius) && (posInUnit.y < radius) {
			multiplier := 1.0
			if (roundMult > 1.0) && (roundMult <= 2) {
				multiplier = 1.0 + (sqrt(2.0)-1.0)*easeOutCirc(roundMult-1.0)
			} else if roundMult < 0 {
				multiplier = sqrt(2) - (sqrt(2.0)-1.0)*easeInSine(-roundMult)
			}
			if posInUnit.y < growUp(posInUnit.x, radius, multiplier) {
				normColor.a = 0.0
			}

		}
	} else {
		// if TotalSize.x >= TotalSize.y {
		radius := TotalSize.y / 2.0

		// Top Left Corner
		if roundMult := ShadedCorners[0]; (roundMult != 0.0) && (posInUnit.x < radius) && (posInUnit.y < radius) {
			multiplier := 1.0
			if (roundMult > 1.0) && (roundMult <= 2) {
				multiplier = 1.0 + (sqrt(2.0)-1.0)*easeOutCirc(roundMult-1.0)
			} else if roundMult < 0 {
				multiplier = sqrt(2) - (sqrt(2.0)-1.0)*easeInSine(-roundMult)
			}
			if posInUnit.x < growLeft(posInUnit.y, radius, multiplier) {
				normColor.a = 0.0
			}
		}
		// Bottom Left Corner
		if roundMult := ShadedCorners[1]; (roundMult != 0.0) && (posInUnit.x < radius) && (posInUnit.y > radius) {
			multiplier := 1.0
			if (roundMult > 1.0) && (roundMult <= 2) {
				multiplier = 1.0 + (sqrt(2.0)-1.0)*easeOutCirc(roundMult-1.0)
			} else if roundMult < 0 {
				multiplier = sqrt(2) - (sqrt(2.0)-1.0)*easeInSine(-roundMult)
			}
			if posInUnit.x < growLeft(posInUnit.y, radius, multiplier) {
				normColor.a = 0.0
			}
		}
		// Bottom Right Corner
		if roundMult := ShadedCorners[2]; (roundMult != 0.0) && (posInUnit.x > TotalSize.x-radius) && (posInUnit.y > radius) {
			multiplier := 1.0
			if (roundMult > 1.0) && (roundMult <= 2) {
				multiplier = 1.0 + (sqrt(2.0)-1.0)*easeOutCirc(roundMult-1.0)
			} else if roundMult < 0 {
				multiplier = sqrt(2) - (sqrt(2.0)-1.0)*easeInSine(-roundMult)
			}
			if posInUnit.x > growRight(posInUnit.y, radius, multiplier) {
				normColor.a = 0.0
			}
		}
		// Top Right Corner
		if roundMult := ShadedCorners[3]; (roundMult != 0.0) && (posInUnit.x > TotalSize.x-radius) && (posInUnit.y < radius) {
			multiplier := 1.0
			if (roundMult > 1.0) && (roundMult <= 2) {
				multiplier = 1.0 + (sqrt(2.0)-1.0)*easeOutCirc(roundMult-1.0)
			} else if roundMult < 0 {
				multiplier = sqrt(2) - (sqrt(2.0)-1.0)*easeInSine(-roundMult)
			}
			if posInUnit.x > growRight(posInUnit.y, radius, multiplier) {
				normColor.a = 0.0
			}
		}
	}

	normColor.rgb *= normColor.a
	return normColor
}

func growUp(x, radius, multiplier float) float {
	return radius - sqrt(pow(radius*multiplier, 2.0)-pow(x-radius, 2.0))
}

func growDown(x, radius, multiplier float) float {
	return TotalSize.y - radius + sqrt(pow(radius*multiplier, 2.0)-pow(x-radius, 2.0))
}

func growLeft(y, radius, multiplier float) float {
	return radius - sqrt(pow(radius*multiplier, 2.0)-pow(y-radius, 2.0))
}

func growRight(y, radius, multiplier float) float {
	return TotalSize.x - radius + sqrt(pow(radius*multiplier, 2.0)-pow(y-radius, 2.0))
}

// https://easings.net/#easeOutCirc
func easeOutCirc(x float) float {
	return sqrt(1.0 - pow(x-1.0, 2.0))
}

// https://easings.net/#easeInSine
func easeInSine(x float) float {
	return 1.0 - cos((x*3.14159)/2.0)
}
