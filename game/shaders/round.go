// Copyright (c) 2022 Anıl Konaç
// This file is licensed under the MIT license.

//go:build ignore

package main

var Color vec4
var ShadedCorners [4]float
var RectSize vec2
var RectPosInUnit vec2
var TotalSize vec2

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	normColor := Color / 0xffff

	posInUnit := vec2(RectSize.x*texCoord.x, RectSize.y*texCoord.y) + RectPosInUnit
	if TotalSize.x <= TotalSize.y {
		radius := TotalSize.x / 2
		roundCenter1 := vec2(radius)
		roundCenter2 := vec2(radius, TotalSize.y-radius)

		distToCenter1 := distance(posInUnit, roundCenter1)
		distToCenter2 := distance(posInUnit, roundCenter2)

		// Top Left Corner
		if (ShadedCorners[0] > 0) && (posInUnit.x < radius) &&
			(posInUnit.y < radius) && (distToCenter1 > radius) {
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
			(posInUnit.y < radius) && (distToCenter1 > radius) {
			normColor.a = 0
		}
	} else {
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
