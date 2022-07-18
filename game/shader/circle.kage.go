//go:build ignore

package main

var (
	Radius float
	Color vec4
)

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	clr := Color
	center := vec2(Radius)
	if distance(texCoord, center) >= Radius {
		clr.a = 0.0
	}

	clr.rgb *= clr.a
	return clr
}