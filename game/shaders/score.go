//go:build ignore

package main

var Alpha float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	imgColor := imageSrc0At(texCoord)
	normColor := color / 0xffff
	var normAlpha float
	if imgColor == normColor {
		normAlpha = 0.0
	} else {
		normAlpha = Alpha / 0xff
	}
	imgColor.rgb *= normAlpha
	imgColor.a = normAlpha
	return imgColor
}
