//go:build ignore

package main

var Alpha float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	imgColor := imageSrc0At(texCoord)

	// Get font drawing information from red color
	redIntensity := imgColor.r

	fontColor := color
	// Set alpha of font color to uniform variable
	fontColor *= Alpha

	// Interpolate between the font color and the full transparent color according to the red channel.
	newColor := redIntensity*fontColor + (1.0-redIntensity)*vec4(0.0)

	return newColor
}
