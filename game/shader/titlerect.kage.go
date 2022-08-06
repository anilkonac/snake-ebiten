//go:build ignore

package main

var RadiusTex vec2

func Fragment(pos vec4, texCoord vec2, col vec4) vec4 {

	// Round corners
	origin, size := imageSrcRegionOnTexture()
	centerTopLeft := origin + RadiusTex
	centerTopRight := vec2(origin.x+size.x-RadiusTex.x, origin.y+RadiusTex.y)
	centerBottomLeft := vec2(origin.x+RadiusTex.x, origin.y+size.y-RadiusTex.y)
	centerBottomRight := vec2(origin.x+size.x-RadiusTex.x, origin.y+size.y-RadiusTex.y)

	if ((texCoord.x < centerTopLeft.x) && (texCoord.y < centerTopLeft.y) && distance(texCoord, centerTopLeft) > RadiusTex.x) ||
		((texCoord.x > centerTopRight.x) && (texCoord.y < centerTopRight.y) && (distance(texCoord, centerTopRight) > RadiusTex.x)) ||
		((texCoord.x < centerBottomLeft.x) && (texCoord.y > centerBottomLeft.y) && (distance(texCoord, centerBottomLeft) > RadiusTex.x)) ||
		((texCoord.x > centerBottomRight.x) && (texCoord.y > centerBottomRight.y) && (distance(texCoord, centerBottomRight) > RadiusTex.x)) {
		return vec4(0.0)
	}

	// Draw rectangle image with texts
	return imageSrc0UnsafeAt(texCoord)

}
