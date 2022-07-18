//go:build ignore

package main

var Alpha float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	imgColor := imageSrc0UnsafeAt(texCoord)
	return imgColor.r * color * Alpha
}
