//go:build ignore

package main

var RawAlpha float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	imgColor := imageSrc0At(texCoord)
	redIntensity := imgColor.r // Get font drawing information from red color

	fontColor := color / 0xffff   // Normalize the score color between 0.0 and 1.0
	fontColor.a = RawAlpha / 0xff // Normalize alpha between 0.0 and 1.0
	fontColor.rgb *= fontColor.a  // Apply alpha to fontColor

	// Interpolate between the font color and the full transparent color according to the red channel.
	newColor := redIntensity*fontColor + (1-redIntensity)*vec4(0.0)

	return newColor
}
