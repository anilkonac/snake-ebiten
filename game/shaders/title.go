//go:build ignore

package main

var ShowKeyPrompt float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	alpha := color.a

	var imgColor vec4
	if ShowKeyPrompt == 1.0 {
		imgColor = imageSrc1UnsafeAt(texCoord)
	} else {
		imgColor = imageSrc0UnsafeAt(texCoord)
	}
	imgColor *= alpha

	return imgColor
}
