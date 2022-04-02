//go:build ignore

package main

var (
	Radius     float
	IsVertical float
	Dimension  vec2
)

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	normColor := color / 0xffff
	alpha := 0.0

	headCenter1 := vec2(Radius, Radius)
	if Dimension.x == Dimension.y {
		if distance(texCoord, headCenter1) <= Radius {
			alpha = normColor.a
		}
	} else if IsVertical > 0.0 {
		headCenter2 := vec2(Radius, Dimension.y-Radius)

		if (texCoord.y < headCenter1.y) && (distance(texCoord, headCenter1) <= Radius) {
			alpha = normColor.a
		} else if (texCoord.y >= headCenter1.y) && (texCoord.y <= headCenter2.y) {
			alpha = normColor.a
		} else if (texCoord.y > headCenter2.y) && (distance(texCoord, headCenter2) <= Radius) {
			alpha = normColor.a
		}
	} else {
		headCenter2 := vec2(Dimension.x-Radius, Radius)

		if (texCoord.x < headCenter1.x) && (distance(texCoord, headCenter1) <= Radius) {
			alpha = normColor.a
		} else if (texCoord.x >= headCenter1.x) && (texCoord.x <= headCenter2.x) {
			alpha = normColor.a
		} else if (texCoord.x > headCenter2.x) && (distance(texCoord, headCenter2) <= Radius) {
			alpha = normColor.a
		}
	}

	normColor.a = alpha
	normColor.rgb *= normColor.a
	return normColor
}
