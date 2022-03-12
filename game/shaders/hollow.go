// Copyright (c) 2022 Anıl Konaç
// This file is licensed under the MIT license.

//go:build ignore

package main

var Color vec4
var RectSize vec2
var RectPosInUnit vec2
var TotalSize vec2
var Thickness float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	normColor := Color / 0xffff

	posInUnit := vec2(RectSize.x*texCoord.x, RectSize.y*texCoord.y) + RectPosInUnit
	if (posInUnit.x >= Thickness) && (posInUnit.x <= TotalSize.x-Thickness) &&
		(posInUnit.y >= Thickness) && (posInUnit.y <= TotalSize.y-Thickness) {
		normColor.a = 0
	}

	normColor.rgb *= normColor.a
	return normColor
}
