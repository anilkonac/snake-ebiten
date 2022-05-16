// Code generated by file2byteslice. DO NOT EDIT.

package shaders

var SnakeHead = []byte("//go:build ignore\r\n\r\npackage main\r\n\r\nvar (\r\n\tRadius      float\r\n\tRadiusMouth float\r\n\tDirection   float\r\n\tProxToFood  float\r\n\tSize        vec2\r\n)\r\n\r\nfunc Fragment(position vec4, texCoord vec2, color vec4) vec4 {\r\n\tclr := color\r\n\r\n\theadCenter1 := vec2(Radius, Radius)\r\n\tif Direction <= 1.0 { // direction is vertical\r\n\t\theadCenter2 := vec2(Radius, Size.y-Radius)\r\n\r\n\t\t// Round the unit\r\n\t\tif (texCoord.y < headCenter1.y) && (distance(texCoord, headCenter1) > Radius) {\r\n\t\t\tclr.a = 0.0\r\n\t\t} else if (texCoord.y > headCenter2.y) && (distance(texCoord, headCenter2) > Radius) {\r\n\t\t\tclr.a = 0.0\r\n\t\t}\r\n\r\n\t\t// Draw mouth\r\n\t\tif isMouthVertical(texCoord) {\r\n\t\t\tclr.a = 0.0\r\n\t\t}\r\n\r\n\t} else { // direction is horizontal\r\n\t\theadCenter2 := vec2(Size.x-Radius, Radius)\r\n\r\n\t\t// Round the unit\r\n\t\tif (texCoord.x < headCenter1.x) && (distance(texCoord, headCenter1) > Radius) {\r\n\t\t\tclr.a = 0.0\r\n\t\t} else if (texCoord.x > headCenter2.x) && (distance(texCoord, headCenter2) > Radius) {\r\n\t\t\tclr.a = 0.0\r\n\t\t}\r\n\r\n\t\t// Draw mouth\r\n\t\tif isMouthHorizontal(texCoord) {\r\n\t\t\tclr.a = 0.0\r\n\t\t}\r\n\r\n\t}\r\n\r\n\tclr.rgb *= clr.a\r\n\treturn clr\r\n}\r\n\r\nfunc isMouthVertical(texCoord vec2) bool {\r\n\t// If the food is far away, don't bother checking if tex is mouth\r\n\tif ProxToFood <= 0.0 {\r\n\t\treturn false\r\n\t}\r\n\r\n\t// Calculate mouth center\r\n\tvar mouthCenter vec2\r\n\tif Direction == 0.0 { // up\r\n\t\tmouthCenter = vec2(Radius, 0.0)\r\n\t} else { // down\r\n\t\tmouthCenter = vec2(Radius, Size.y)\r\n\t}\r\n\r\n\t// Check if the position is in the mouth\r\n\tif distance(texCoord, mouthCenter) < RadiusMouth*easeOutCubic(ProxToFood) {\r\n\t\treturn true\r\n\t}\r\n\treturn false\r\n}\r\n\r\nfunc isMouthHorizontal(texCoord vec2) bool {\r\n\t// If the food is far away, don't bother checking if tex is mouth\r\n\tif ProxToFood <= 0.0 {\r\n\t\treturn false\r\n\t}\r\n\r\n\t// Calculate mouth center\r\n\tvar mouthCenter vec2\r\n\tif Direction == 2.0 { // left\r\n\t\tmouthCenter = vec2(0.0, Radius)\r\n\t} else { // right\r\n\t\tmouthCenter = vec2(Size.x, Radius)\r\n\t}\r\n\r\n\t// Check if the position is in the mouth\r\n\tif distance(texCoord, mouthCenter) < RadiusMouth*easeOutCubic(ProxToFood) {\r\n\t\treturn true\r\n\t}\r\n\treturn false\r\n}\r\n\r\n// https://easings.net/#easeOutCubic\r\nfunc easeOutCubic(x float) float {\r\n\txMin := 1.0 - x\r\n\treturn 1.0 - xMin*xMin*xMin\r\n}\r\n")
