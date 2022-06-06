//go:build ignore

package main

var (
	Radius    float
	Direction float
	Size      vec2
)

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	clr := color
	isVertical := (Direction <= 1.0)

	headCenter1 := vec2(Radius, Radius)
	if Size.x == Size.y {
		if distance(texCoord, headCenter1) > Radius {
			clr.a = 0.0
		}
	} else if isVertical {
		headCenter2 := vec2(Radius, Size.y-Radius)

		if (texCoord.y < headCenter1.y) && (distance(texCoord, headCenter1) > Radius) {
			clr.a = 0.0
		} else if (texCoord.y > headCenter2.y) && (distance(texCoord, headCenter2) > Radius) {
			clr.a = 0.0
		}
	} else {
		headCenter2 := vec2(Size.x-Radius, Radius)

		if (texCoord.x < headCenter1.x) && (distance(texCoord, headCenter1) > Radius) {
			clr.a = 0.0
		} else if (texCoord.x > headCenter2.x) && (distance(texCoord, headCenter2) > Radius) {
			clr.a = 0.0
		}
	}

	clr.rgb *= clr.a
	return clr
}