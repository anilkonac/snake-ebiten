// Code generated by file2byteslice. DO NOT EDIT.

package shaders

var Round = []byte("// snake-ebiten\n// Copyright (C) 2022 Anıl Konaç\n\n// This program is free software: you can redistribute it and/or modify\n// it under the terms of the GNU General Public License as published by\n// the Free Software Foundation, either version 3 of the License, or\n// (at your option) any later version.\n\n// This program is distributed in the hope that it will be useful,\n// but WITHOUT ANY WARRANTY; without even the implied warranty of\n// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\n// GNU General Public License for more details.\n\n// You should have received a copy of the GNU General Public License\n// along with this program.  If not, see <https://www.gnu.org/licenses/>.\n\n//go:build ignore\n\npackage main\n\nvar Color vec4\nvar ShadedCorners [4]float\nvar RectSize vec2\nvar RectPosInUnit vec2\nvar TotalSize vec2\nvar IsVertical float\n\nfunc Fragment(position vec4, texCoord vec2, color vec4) vec4 {\n\tnormColor := Color / 0xffff\n\n\tposInUnit := vec2(RectSize.x*texCoord.x, RectSize.y*texCoord.y) + RectPosInUnit\n\tif IsVertical > 0 {\n\t\tradius := TotalSize.x / 2\n\n\t\t// Top Left Corner\n\t\tif roundMult := ShadedCorners[0]; (roundMult > 0) && (posInUnit.x < radius) &&\n\t\t\t(posInUnit.y < radius) && (posInUnit.y < growUp(posInUnit.x, radius)) {\n\t\t\tnormColor.a -= roundMult\n\t\t}\n\t\t// Bottom Left Corner\n\t\tif roundMult := ShadedCorners[1]; (roundMult > 0) && (posInUnit.x < radius) &&\n\t\t\t(posInUnit.y > (TotalSize.y - radius)) && (posInUnit.y > growDown(posInUnit.x, radius)) {\n\t\t\tnormColor.a -= roundMult\n\t\t}\n\t\t// Bottom Right Corner\n\t\tif roundMult := ShadedCorners[2]; (roundMult > 0) && (posInUnit.x > radius) &&\n\t\t\t(posInUnit.y > (TotalSize.y - radius)) && (posInUnit.y > growDown(posInUnit.x, radius)) {\n\t\t\tnormColor.a -= roundMult\n\t\t}\n\t\t// Top Right Corner\n\t\tif roundMult := ShadedCorners[3]; (roundMult > 0) && (posInUnit.x > radius) &&\n\t\t\t(posInUnit.y < radius) && (posInUnit.y < growUp(posInUnit.x, radius)) {\n\t\t\tnormColor.a -= roundMult\n\t\t}\n\t} else {\n\t\t// if TotalSize.x >= TotalSize.y {\n\t\tradius := TotalSize.y / 2\n\n\t\t// Top Left Corner\n\t\tif roundMult := ShadedCorners[0]; (roundMult > 0) && (posInUnit.x < radius) &&\n\t\t\t(posInUnit.y < radius) && (posInUnit.x < growLeft(posInUnit.y, radius)) {\n\t\t\tnormColor.a -= roundMult\n\t\t}\n\t\t// Bottom Left Corner\n\t\tif roundMult := ShadedCorners[1]; (roundMult > 0) && (posInUnit.x < radius) &&\n\t\t\t(posInUnit.y > radius) && (posInUnit.x < growLeft(posInUnit.y, radius)) {\n\t\t\tnormColor.a -= roundMult\n\t\t}\n\t\t// Bottom Right Corner\n\t\tif roundMult := ShadedCorners[2]; (roundMult > 0) && (posInUnit.x > TotalSize.x-radius) &&\n\t\t\t(posInUnit.y > radius) && (posInUnit.x > growRight(posInUnit.y, radius)) {\n\t\t\tnormColor.a -= roundMult\n\t\t}\n\t\t// Top Right Corner\n\t\tif roundMult := ShadedCorners[3]; (roundMult > 0) && (posInUnit.x > TotalSize.x-radius) &&\n\t\t\t(posInUnit.y < radius) && (posInUnit.x > growRight(posInUnit.y, radius)) {\n\t\t\tnormColor.a -= roundMult\n\t\t}\n\t}\n\n\tnormColor.rgb *= normColor.a\n\treturn normColor\n}\n\nfunc growUp(x, radius float) float {\n\t// Interpolate between square function and semicircle function\n\theightMultip := clamp(TotalSize.y/radius, 0.0, 1.0)\n\n\tsquare := squareGrowthReverse(x, radius, heightMultip)\n\tsemicircle := semicircleReverse(x, radius, heightMultip)\n\n\ttransition := pow(heightMultip, 5.0) // easeInQuint\n\treturn transition*semicircle + (1.0-transition)*square\n}\n\nfunc growDown(x, radius float) float {\n\t// Interpolate between square function and semicircle function\n\theightMultip := clamp(TotalSize.y/radius, 0.0, 1.0)\n\n\tsquare := TotalSize.y - heightMultip*radius + squareGrowth(x, radius, heightMultip)\n\tsemicircle := TotalSize.y - radius + semicircle(x, radius, heightMultip)\n\n\ttransition := pow(heightMultip, 5.0) // easeInQuint\n\treturn transition*semicircle + (1.0-transition)*square\n}\n\nfunc growLeft(y, radius float) float {\n\t// Interpolate between square function and semicircle function\n\twidthMultip := clamp(TotalSize.x/radius, 0.0, 1.0)\n\n\tsquare := squareGrowthReverse(y, radius, widthMultip)\n\tsemicircle := semicircleReverse(y, radius, widthMultip)\n\n\ttransition := pow(widthMultip, 5.0) // easeInQuint\n\treturn transition*semicircle + (1.0-transition)*square\n}\n\nfunc growRight(y, radius float) float {\n\t// Interpolate between square function and semicircle function\n\twidthMultip := clamp(TotalSize.x/radius, 0.0, 1.0)\n\n\tsquare := TotalSize.x - widthMultip*radius + squareGrowth(y, radius, widthMultip)\n\tsemicircle := TotalSize.x - radius + semicircle(y, radius, widthMultip)\n\n\ttransition := pow(widthMultip, 5.0) // easeInQuint\n\treturn transition*semicircle + (1.0-transition)*square\n}\n\nfunc squareGrowth(x, radius, multiplier float) float {\n\treturn multiplier * (radius - pow(x-radius, 2.0)/radius)\n}\n\nfunc squareGrowthReverse(x, radius, multiplier float) float {\n\treturn multiplier * pow(x-radius, 2.0) / radius\n}\n\nfunc semicircle(x, radius, multiplier float) float {\n\treturn sqrt(pow(radius, 2.0) - pow(x-radius, 2.0))\n}\n\nfunc semicircleReverse(x, radius, multiplier float) float {\n\treturn radius - semicircle(x, radius, multiplier)\n}\n")
