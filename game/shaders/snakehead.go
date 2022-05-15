//go:build ignore

package main

var (
	Radius      float
	RadiusMouth float
	Direction   float
	DistToFood  float
	Size        vec2
)

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	clr := color

	headCenter1 := vec2(Radius, Radius)
	if Direction <= 1.0 { // direction is vertical
		headCenter2 := vec2(Radius, Size.y-Radius)

		if DistToFood < 1.0 {
			// Calculate mouth center
			var mouthCenter vec2
			if Direction == 0.0 { // up
				mouthCenter = vec2(Radius, 0.0)
			} else { // down
				mouthCenter = vec2(Radius, Size.y)
			}

			// Draw mouth
			if distance(texCoord, mouthCenter) < RadiusMouth*(1.0-DistToFood) {
				clr.a = 0.0
			}
		}

		// Round the unit
		if (texCoord.y < headCenter1.y) && (distance(texCoord, headCenter1) > Radius) {
			clr.a = 0.0
		} else if (texCoord.y > headCenter2.y) && (distance(texCoord, headCenter2) > Radius) {
			clr.a = 0.0
		}

	} else { // direction is horizontal
		headCenter2 := vec2(Size.x-Radius, Radius)

		if DistToFood < 1.0 {
			// Calculate mouth center
			var mouthCenter vec2
			if Direction == 2.0 { // left
				mouthCenter = vec2(0.0, Radius)
			} else { // right
				mouthCenter = vec2(Size.x, Radius)
			}

			// Draw mouth
			if distance(texCoord, mouthCenter) < RadiusMouth*(1.0-DistToFood) {
				clr.a = 0.0
			}
		}

		// Round the unit
		if (texCoord.x < headCenter1.x) && (distance(texCoord, headCenter1) > Radius) {
			clr.a = 0.0
		} else if (texCoord.x > headCenter2.x) && (distance(texCoord, headCenter2) > Radius) {
			clr.a = 0.0
		}

	}

	clr.rgb *= clr.a
	return clr
}