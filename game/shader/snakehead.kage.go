//go:build ignore

package main

var (
	RadiusMouth float
	Direction   float
	ProxToFood  float
	Size        vec2
)

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	clr := color

	if Direction <= 1.0 { // direction is vertical
		radius := Size.x / 2.0
		headCenter1 := vec2(radius)
		headCenter2 := vec2(radius, Size.y-radius)

		// Round the unit
		if (texCoord.y < headCenter1.y) && (distance(texCoord, headCenter1) > radius) {
			clr.a = 0.0
		} else if (texCoord.y > headCenter2.y) && (distance(texCoord, headCenter2) > radius) {
			clr.a = 0.0
		}

		// Draw mouth
		if isMouthVertical(texCoord, radius) {
			clr.a = 0.0
		}

	} else { // direction is horizontal
		radius := Size.y / 2.0
		headCenter1 := vec2(radius)
		headCenter2 := vec2(Size.x-radius, radius)

		// Round the unit
		if (texCoord.x < headCenter1.x) && (distance(texCoord, headCenter1) > radius) {
			clr.a = 0.0
		} else if (texCoord.x > headCenter2.x) && (distance(texCoord, headCenter2) > radius) {
			clr.a = 0.0
		}

		// Draw mouth
		if isMouthHorizontal(texCoord, radius) {
			clr.a = 0.0
		}

	}

	clr.rgb *= clr.a
	return clr
}

func isMouthVertical(texCoord vec2, radius float) bool {
	// If the food is far away, don't bother checking if tex is mouth
	if ProxToFood <= 0.0 {
		return false
	}

	// Calculate mouth center
	var mouthCenter vec2
	if Direction == 0.0 { // up
		mouthCenter = vec2(radius, 0.0)
	} else { // down
		mouthCenter = vec2(radius, Size.y)
	}

	// Check if the position is in the mouth
	if distance(texCoord, mouthCenter) < RadiusMouth*easeOutCubic(ProxToFood) {
		return true
	}
	return false
}

func isMouthHorizontal(texCoord vec2, radius float) bool {
	// If the food is far away, don't bother checking if tex is mouth
	if ProxToFood <= 0.0 {
		return false
	}

	// Calculate mouth center
	var mouthCenter vec2
	if Direction == 2.0 { // left
		mouthCenter = vec2(0.0, radius)
	} else { // right
		mouthCenter = vec2(Size.x, radius)
	}

	// Check if the position is in the mouth
	if distance(texCoord, mouthCenter) < RadiusMouth*easeOutCubic(ProxToFood) {
		return true
	}
	return false
}

// https://easings.net/#easeOutCubic
func easeOutCubic(x float) float {
	xMin := 1.0 - x
	return 1.0 - xMin*xMin*xMin
}