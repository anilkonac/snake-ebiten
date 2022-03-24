//go:build ignore

package main

var Color vec4

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	normColor := Color / 0xffff
	normColor.rgb *= normColor.a
	return normColor
}
