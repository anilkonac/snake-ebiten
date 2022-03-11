// Copyright (c) 2022 Anıl Konaç
// This file is licensed under the MIT license.

//go:build ignore

package main

var Color vec4
var Thickness float
var Width float
var Height float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	normColor := Color / 0xffff

	realX := Width * texCoord.x
	realY := Height * texCoord.y
	if (realX >= Thickness) && (realX <= Width-Thickness) &&
		(realY >= Thickness) && (realY <= Height-Thickness) {
		normColor.a = 0
	}

	normColor.xyz *= normColor.a
	return normColor
}
