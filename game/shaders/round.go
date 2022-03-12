// Copyright (c) 2022 Anıl Konaç
// This file is licensed under the MIT license.

//go:build ignore

package main

var Color vec4
var Width float
var Height float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	normColor := Color / 0xffff

	realCoord := vec2(Width*texCoord.x, Height*texCoord.y)
	if Width <= Height {
		radius := Width / 2
		roundCenter1 := vec2(radius)
		roundCenter2 := vec2(radius, Height-radius)

		distToCenter1 := distance(realCoord, roundCenter1)
		distToCenter2 := distance(realCoord, roundCenter2)
		if ((realCoord.y < radius) && (distToCenter1 > radius)) ||
			((realCoord.y > (Height - radius)) && (distToCenter2 > radius)) {
			normColor.a = 0
		}
	} else {
		radius := Height / 2
		roundCenter1 := vec2(radius)
		roundCenter2 := vec2(Width-radius, radius)

		distToCenter1 := distance(realCoord, roundCenter1)
		distToCenter2 := distance(realCoord, roundCenter2)
		if ((realCoord.x < radius) && (distToCenter1 > radius)) ||
			((realCoord.x > (Width - radius)) && (distToCenter2 > radius)) {
			normColor.a = 0
		}
	}

	normColor.xyz *= normColor.a
	return normColor
}
