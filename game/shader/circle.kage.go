//go:build ignore

package main

var (
	Radius float
)

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	center := vec2(Radius)
	if distance(texCoord, center) > Radius {
		return vec4(0.0)
	}

	return vec4(1.0)
}
