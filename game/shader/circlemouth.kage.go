//go:build ignore

package main

var (
	Radius      float
	RadiusMouth float
	Direction   float
	ProxToFood  float
)

func Fragment(pos vec4, tex vec2, col vec4) vec4 {
	clr := col

	// Round the circle
	center := vec2(Radius)
	if distance(tex, center) > Radius {
		clr.a = 0.0
	}

	// Draw mouth
	// ----------
	if ProxToFood <= 0.0 {
		clr.rgb *= clr.a
		return clr
	}

	// Specify mouth center
	var mouthCenter vec2
	if Direction == 0.0 { // up
		mouthCenter = vec2(Radius, 0.0)
	} else if Direction == 1.0 { // down
		mouthCenter = vec2(Radius, 2*Radius)
	} else if Direction == 2.0 { // left
		mouthCenter = vec2(0.0, Radius)
	} else if Direction == 3.0 { // right
		mouthCenter = vec2(2*Radius, Radius)
	}

	// Check if the position is in the mouth
	if distance(tex, mouthCenter) < RadiusMouth*easeOutCubic(ProxToFood) {
		clr.a = 0.0
	}

	clr.rgb *= clr.a
	return clr
}

// https://easings.net/#easeOutCubic
func easeOutCubic(x float) float {
	xMin := 1.0 - x
	return 1.0 - xMin*xMin*xMin
}
