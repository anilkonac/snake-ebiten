//go:build ignore

package main

var (
	Radius     float
	Direction  float
	DistToFood float
	Size       vec2
)

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	clr := color
	mouthRadius := Radius * 1.10

	headCenter1 := vec2(Radius, Radius)
	var headCenter2 vec2
	if Direction <= 1.0 { // if direction is vertical
		headCenter2 = vec2(Radius, Size.y-Radius)
	} else {
		headCenter2 = vec2(Size.x-Radius, Radius)
	}

	if Direction == 0.0 { // up
		mouthCenter := vec2(Radius, -DistToFood*mouthRadius)
		if (texCoord.y < headCenter1.y) && ((distance(texCoord, headCenter1) > Radius) || (distance(texCoord, mouthCenter) < mouthRadius)) {
			clr.a = 0.0
		} else if (texCoord.y > headCenter2.y) && (distance(texCoord, headCenter2) > Radius) {
			clr.a = 0.0
		}
	} else if Direction == 1.0 { // down
		mouthCenter := vec2(Radius, Size.y+DistToFood*mouthRadius)
		if (texCoord.y < headCenter1.y) && (distance(texCoord, headCenter1) > Radius) {
			clr.a = 0.0
		} else if (texCoord.y > headCenter2.y) && ((distance(texCoord, headCenter2) > Radius) || (distance(texCoord, mouthCenter) < mouthRadius)) {
			clr.a = 0.0
		}
	} else if Direction == 2.0 { // left
		mouthCenter := vec2(-DistToFood*mouthRadius, Radius)
		if (texCoord.x < headCenter1.x) && ((distance(texCoord, headCenter1) > Radius) || (distance(texCoord, mouthCenter) < mouthRadius)) {
			clr.a = 0.0
		} else if (texCoord.x > headCenter2.x) && (distance(texCoord, headCenter2) > Radius) {
			clr.a = 0.0
		}
	} else { // right
		mouthCenter := vec2(Size.x+DistToFood*mouthRadius, Radius)
		if (texCoord.x < headCenter1.x) && (distance(texCoord, headCenter1) > Radius) {
			clr.a = 0.0
		} else if (texCoord.x > headCenter2.x) && ((distance(texCoord, headCenter2) > Radius) || (distance(texCoord, mouthCenter) < mouthRadius)) {
			clr.a = 0.0
		}
	}

	clr.rgb *= clr.a
	return clr
}
