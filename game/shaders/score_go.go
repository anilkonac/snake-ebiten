// Code generated by file2byteslice. DO NOT EDIT.

package shaders

var Score = []byte("//go:build ignore\n\npackage main\n\nvar Alpha float\n\nfunc Fragment(position vec4, texCoord vec2, color vec4) vec4 {\n\timgColor := imageSrc0At(texCoord)\n\n\t// Get font drawing information from red color\n\tredIntensity := imgColor.r\n\n\tfontColor := color\n\t// Set alpha of font color to uniform variable\n\tfontColor *= Alpha\n\n\t// Interpolate between font color and full transparent color according to the red channel.\n\treturn redIntensity*fontColor + (1.0-redIntensity)*vec4(0.0)\n}\n")
