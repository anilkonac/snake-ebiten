//go:build ignore

package main

var (
	CornerRadius vec2
	Size         vec2
)

func Fragment(pos vec4, texCoord vec2, col vec4) vec4 {

	// Round the corners of the rectangle image
	centerTopLeft := CornerRadius
	centerTopRight := vec2(Size.x-CornerRadius.x, CornerRadius.y)
	centerBottomLeft := vec2(CornerRadius.x, Size.y-CornerRadius.y)
	centerBottomRight := vec2(Size.x-CornerRadius.x, Size.y-CornerRadius.y)

	if ((texCoord.x < centerTopLeft.x) && (texCoord.y < centerTopLeft.y) && distance(texCoord, centerTopLeft) > CornerRadius.x) ||
		((texCoord.x > centerTopRight.x) && (texCoord.y < centerTopRight.y) && (distance(texCoord, centerTopRight) > CornerRadius.x)) ||
		((texCoord.x < centerBottomLeft.x) && (texCoord.y > centerBottomLeft.y) && (distance(texCoord, centerBottomLeft) > CornerRadius.x)) ||
		((texCoord.x > centerBottomRight.x) && (texCoord.y > centerBottomRight.y) && (distance(texCoord, centerBottomRight) > CornerRadius.x)) {
		return vec4(0.0)
	}

	return vec4(1.0)

}
