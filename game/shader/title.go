//go:build ignore

package main

var (
	ShowKeyPrompt float
	Alpha         float
	RadiusTex     vec2
)

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {

	// Get the color value of the current texel
	var imgColor vec4
	if ShowKeyPrompt == 1.0 {
		imgColor = imageSrc1UnsafeAt(texCoord)
	} else {
		imgColor = imageSrc0UnsafeAt(texCoord)
	}

	// Round the corners of the rectangle image
	origin, size := imageSrcRegionOnTexture()
	centerTopLeft := origin + RadiusTex
	centerTopRight := vec2(origin.x+size.x-RadiusTex.x, origin.y+RadiusTex.y)
	centerBottomLeft := vec2(origin.x+RadiusTex.x, origin.y+size.y-RadiusTex.y)
	centerBottomRight := vec2(origin.x+size.x-RadiusTex.x, origin.y+size.y-RadiusTex.y)

	alpha := Alpha
	if (texCoord.x < centerTopLeft.x) && (texCoord.y < centerTopLeft.y) && distance(texCoord, centerTopLeft) > RadiusTex.x {
		alpha = 0.0
	} else if (texCoord.x > centerTopRight.x) && (texCoord.y < centerTopRight.y) && (distance(texCoord, centerTopRight) > RadiusTex.x) {
		alpha = 0.0
	} else if (texCoord.x < centerBottomLeft.x) && (texCoord.y > centerBottomLeft.y) && (distance(texCoord, centerBottomLeft) > RadiusTex.x) {
		alpha = 0.0
	} else if (texCoord.x > centerBottomRight.x) && (texCoord.y > centerBottomRight.y) && (distance(texCoord, centerBottomRight) > RadiusTex.x) {
		alpha = 0.0
	}

	return imgColor * alpha
}
