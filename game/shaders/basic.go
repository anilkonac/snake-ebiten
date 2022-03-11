//go:build ignore

package main

var Color vec4
var TimeElapsed float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	a := Color.a / 0xffff
	r := (Color.r / 0xffff) * a
	g := (Color.g / 0xffff) * a
	b := (Color.b / 0xffff) * a
	return vec4(r, g, b, a)
}
