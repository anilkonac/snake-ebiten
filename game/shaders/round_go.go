// Code generated by file2byteslice. DO NOT EDIT.

package shaders

var Round = []byte("// snake-ebiten\r\n// Copyright (C) 2022 Anıl Konaç\r\n\r\n// This program is free software: you can redistribute it and/or modify\r\n// it under the terms of the GNU General Public License as published by\r\n// the Free Software Foundation, either version 3 of the License, or\r\n// (at your option) any later version.\r\n\r\n// This program is distributed in the hope that it will be useful,\r\n// but WITHOUT ANY WARRANTY; without even the implied warranty of\r\n// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\r\n// GNU General Public License for more details.\r\n\r\n// You should have received a copy of the GNU General Public License\r\n// along with this program.  If not, see <https://www.gnu.org/licenses/>.\r\n\r\n//go:build ignore\r\n\r\npackage main\r\n\r\nvar Color vec4\r\nvar ShadedCorners [4]float\r\nvar RectSize vec2\r\nvar RectPosInUnit vec2\r\nvar TotalSize vec2\r\nvar IsVertical float\r\n\r\nfunc Fragment(position vec4, texCoord vec2, color vec4) vec4 {\r\n\tnormColor := Color / 0xffff\r\n\r\n\tposInUnit := vec2(RectSize.x*texCoord.x, RectSize.y*texCoord.y) + RectPosInUnit\r\n\tif IsVertical > 0 {\r\n\t\tradius := TotalSize.x / 2\r\n\r\n\t\t// Top Left Corner\r\n\t\tif (ShadedCorners[0] > 0) && (posInUnit.x < radius) &&\r\n\t\t\t(posInUnit.y < radius) && (posInUnit.y < growUp(posInUnit.x, radius)) {\r\n\t\t\tnormColor.a = 0\r\n\t\t}\r\n\t\t// Bottom Left Corner\r\n\t\tif (ShadedCorners[1] > 0) && (posInUnit.x < radius) &&\r\n\t\t\t(posInUnit.y > (TotalSize.y - radius)) && (posInUnit.y > growDown(posInUnit.x, radius)) {\r\n\t\t\tnormColor.a = 0\r\n\t\t}\r\n\t\t// Bottom Right Corner\r\n\t\tif (ShadedCorners[2] > 0) && (posInUnit.x > radius) &&\r\n\t\t\t(posInUnit.y > (TotalSize.y - radius)) && (posInUnit.y > growDown(posInUnit.x, radius)) {\r\n\t\t\tnormColor.a = 0\r\n\t\t}\r\n\t\t// Top Right Corner\r\n\t\tif (ShadedCorners[3] > 0) && (posInUnit.x > radius) &&\r\n\t\t\t(posInUnit.y < radius) && (posInUnit.y < growUp(posInUnit.x, radius)) {\r\n\t\t\tnormColor.a = 0\r\n\t\t}\r\n\r\n\t} else {\r\n\t\t// if TotalSize.x >= TotalSize.y {\r\n\t\tradius := TotalSize.y / 2\r\n\r\n\t\t// Top Left Corner\r\n\t\tif (ShadedCorners[0] > 0) && (posInUnit.x < radius) &&\r\n\t\t\t(posInUnit.y < radius) && (posInUnit.x < growLeft(posInUnit.y, radius)) {\r\n\t\t\tnormColor.a = 0\r\n\t\t}\r\n\t\t// Bottom Left Corner\r\n\t\tif (ShadedCorners[1] > 0) && (posInUnit.x < radius) &&\r\n\t\t\t(posInUnit.y > radius) && (posInUnit.x < growLeft(posInUnit.y, radius)) {\r\n\t\t\tnormColor.a = 0\r\n\t\t}\r\n\t\t// Bottom Right Corner\r\n\t\tif (ShadedCorners[2] > 0) && (posInUnit.x > TotalSize.x-radius) &&\r\n\t\t\t(posInUnit.y > radius) && (posInUnit.x > growRight(posInUnit.y, radius)) {\r\n\t\t\tnormColor.a = 0\r\n\t\t}\r\n\t\t// Top Right Corner\r\n\t\tif (ShadedCorners[3] > 0) && (posInUnit.x > TotalSize.x-radius) &&\r\n\t\t\t(posInUnit.y < radius) && (posInUnit.x > growRight(posInUnit.y, radius)) {\r\n\t\t\tnormColor.a = 0\r\n\t\t}\r\n\t}\r\n\r\n\tnormColor.rgb *= normColor.a\r\n\treturn normColor\r\n}\r\n\r\nfunc growUp(x, radius float) float {\r\n\t// Linear interpolation between square function and semicircle function\r\n\theightMultip := clamp(TotalSize.y/radius, 0.0, 1.0)\r\n\tsquare := squareGrowthReverse(x, radius, heightMultip)\r\n\tsemicircle := semicircleReverse(x, radius, heightMultip)\r\n\r\n\treturn heightMultip*semicircle + (1.0-heightMultip)*square\r\n}\r\n\r\nfunc growDown(x, radius float) float {\r\n\t// Linear interpolation between square function and semicircle function\r\n\theightMultip := clamp(TotalSize.y/radius, 0.0, 1.0)\r\n\tsquare := squareGrowth(x, radius, heightMultip)\r\n\tsemicircle := TotalSize.y - radius + semicircle(x, radius, heightMultip)\r\n\r\n\treturn heightMultip*semicircle + (1.0-heightMultip)*square\r\n}\r\n\r\nfunc growLeft(y, radius float) float {\r\n\t// Linear interpolation between square function and semicircle function\r\n\twidthMultip := clamp(TotalSize.x/radius, 0.0, 1.0)\r\n\tsquare := squareGrowthReverse(y, radius, widthMultip)\r\n\tsemicircle := semicircleReverse(y, radius, widthMultip)\r\n\r\n\treturn widthMultip*semicircle + (1.0-widthMultip)*square\r\n}\r\n\r\nfunc growRight(y, radius float) float {\r\n\t// Linear interpolation between square function and semicircle function\r\n\twidthMultip := clamp(TotalSize.x/radius, 0.0, 1.0)\r\n\tsquare := squareGrowth(y, radius, widthMultip)\r\n\tsemicircle := TotalSize.x - radius + semicircle(y, radius, widthMultip)\r\n\r\n\treturn widthMultip*semicircle + (1.0-widthMultip)*square\r\n}\r\n\r\nfunc squareGrowth(x, radius, multiplier float) float {\r\n\treturn multiplier * (radius - pow(x-radius, 2.0)/radius)\r\n}\r\n\r\nfunc squareGrowthReverse(x, radius, multiplier float) float {\r\n\treturn multiplier * pow(x-radius, 2.0) / radius\r\n}\r\n\r\nfunc semicircle(x, radius, multiplier float) float {\r\n\treturn sqrt(pow(radius, 2.0) - pow(x-radius, 2.0))\r\n}\r\n\r\nfunc semicircleReverse(x, radius, multiplier float) float {\r\n\treturn radius - semicircle(x, radius, multiplier)\r\n}\r\n")
