//go:build ignore

package main

var (
	Direction float
	Size      vec2
)

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	clr := color
	isVertical := (Direction <= 1.0)

	if Size.x == Size.y {
		radius := Size.x / 2.0
		headCenter1 := vec2(radius)
		if distance(texCoord, headCenter1) > radius {
			clr.a = 0.0
		}
	} else if isVertical {
		radius := Size.x / 2.0
		headCenter1 := vec2(radius)
		headCenter2 := vec2(radius, Size.y-radius)

		if (texCoord.y < headCenter1.y) && (distance(texCoord, headCenter1) > radius) {
			clr.a = 0.0
		} else if (texCoord.y > headCenter2.y) && (distance(texCoord, headCenter2) > radius) {
			clr.a = 0.0
		}
	} else {
		radius := Size.y / 2.0
		headCenter1 := vec2(radius)
		headCenter2 := vec2(Size.x-radius, radius)

		if (texCoord.x < headCenter1.x) && (distance(texCoord, headCenter1) > radius) {
			clr.a = 0.0
		} else if (texCoord.x > headCenter2.x) && (distance(texCoord, headCenter2) > radius) {
			clr.a = 0.0
		}
	}

	clr.rgb *= clr.a
	return clr
}
