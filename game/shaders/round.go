// Copyright (c) 2022 Anıl Konaç
// This file is licensed under the MIT license.

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
		// if TotalSize.y >= TotalSize.x {
		radius := TotalSize.x / 2
		// roundCenter1 := vec2(radius)
		roundCenter2 := vec2(radius, TotalSize.y-radius)

		// distToCenter1 := distance(posInUnit, roundCenter1)
		distToCenter2 := distance(posInUnit, roundCenter2)
		multiplier := clamp(TotalSize.y/radius, 0, 1)

		// Top Left Corner
		if (ShadedCorners[0] > 0) && (posInUnit.x < radius) &&
			(posInUnit.y < radius) && (posInUnit.y <= theFunc(posInUnit.x, radius, multiplier)) {
			normColor.a = 0
		}
		// Bottom Left Corner
		if (ShadedCorners[1] > 0) && (posInUnit.x < radius) &&
			(posInUnit.y > (TotalSize.y - radius)) && (distToCenter2 > radius) {
			normColor.a = 0
		}
		// Bottom Right Corner
		if (ShadedCorners[2] > 0) && (posInUnit.x > radius) &&
			(posInUnit.y > (TotalSize.y - radius)) && (distToCenter2 > radius) {
			normColor.a = 0
		}
		// Top Right Corner
		if (ShadedCorners[3] > 0) && (posInUnit.x > radius) &&
			(posInUnit.y < radius) && (posInUnit.y <= theFunc(posInUnit.x, radius, multiplier)) {
			normColor.a = 0
		}
		// }

	} else {
		// if TotalSize.x >= TotalSize.y {
		radius := TotalSize.y / 2
		roundCenter1 := vec2(radius)
		roundCenter2 := vec2(TotalSize.x-radius, radius)

		distToCenter1 := distance(posInUnit, roundCenter1)
		distToCenter2 := distance(posInUnit, roundCenter2)

		// radius *= clamp((TotalSize.y / TotalSize.x), 0, 1)
		// radius = min(radius, TotalSize.x)
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
		// }
	}

	normColor.rgb *= normColor.a
	return normColor
}

func theFunc(x float, radius float, multiplier float) float {
	// if multiplier == 1 {
	// 	return radius - sqrt(radius*radius-(x-radius)*(x-radius))
	// }
	// return multiplier * ( /*radius -*/ pow(x-radius, 2) / radius)

	// Linear interpolation between squared function and circle function
	y0 := multiplier * ( /*radius -*/ pow(x-radius, 2) / radius)
	x0 := 0
	x1 := 1
	y1 := radius - sqrt(radius*radius-(x-radius)*(x-radius))

	return y0 + (multiplier-x0)*(y1-y0)/(x1-x0)
}
