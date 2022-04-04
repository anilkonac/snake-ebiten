//go:build ignore

package main

var (
	Radius     float
	IsVertical float
	Dimension  vec2
)

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	normColor := color / 0xffff

	headCenter1 := vec2(Radius, Radius)
	if Dimension.x == Dimension.y {
		if distance(texCoord, headCenter1) > Radius {
			normColor.a = 0.0
		}
	} else if IsVertical > 0.0 {
		headCenter2 := vec2(Radius, Dimension.y-Radius)

		if (texCoord.y < headCenter1.y) && (distance(texCoord, headCenter1) > Radius) {
			normColor.a = 0.0
		} else if (texCoord.y > headCenter2.y) && (distance(texCoord, headCenter2) > Radius) {
			normColor.a = 0.0
		}
	} else {
		headCenter2 := vec2(Dimension.x-Radius, Radius)

		if (texCoord.x < headCenter1.x) && (distance(texCoord, headCenter1) > Radius) {
			normColor.a = 0.0
		} else if (texCoord.x > headCenter2.x) && (distance(texCoord, headCenter2) > Radius) {
			normColor.a = 0.0
		}
	}

	normColor.rgb *= normColor.a
	return normColor
}
