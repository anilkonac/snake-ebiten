// Code generated by file2byteslice. DO NOT EDIT.

package shaders

var Round = []byte("//go:build ignore\n\npackage main\n\nvar (\n\tRadius     float\n\tIsVertical float\n\tDimension  vec2\n)\n\nfunc Fragment(position vec4, texCoord vec2, color vec4) vec4 {\n\tnormColor := color / 0xffff\n\n\theadCenter1 := vec2(Radius, Radius)\n\tif Dimension.x == Dimension.y {\n\t\tif distance(texCoord, headCenter1) > Radius {\n\t\t\tnormColor.a = 0.0\n\t\t}\n\t} else if IsVertical > 0.0 {\n\t\theadCenter2 := vec2(Radius, Dimension.y-Radius)\n\n\t\tif (texCoord.y < headCenter1.y) && (distance(texCoord, headCenter1) > Radius) {\n\t\t\tnormColor.a = 0.0\n\t\t} else if (texCoord.y > headCenter2.y) && (distance(texCoord, headCenter2) > Radius) {\n\t\t\tnormColor.a = 0.0\n\t\t}\n\t} else {\n\t\theadCenter2 := vec2(Dimension.x-Radius, Radius)\n\n\t\tif (texCoord.x < headCenter1.x) && (distance(texCoord, headCenter1) > Radius) {\n\t\t\tnormColor.a = 0.0\n\t\t} else if (texCoord.x > headCenter2.x) && (distance(texCoord, headCenter2) > Radius) {\n\t\t\tnormColor.a = 0.0\n\t\t}\n\t}\n\n\tnormColor.rgb *= normColor.a\n\treturn normColor\n}\n")
