// Code generated by file2byteslice. DO NOT EDIT.

package shaders

var Round = []byte("// Copyright (c) 2022 Anıl Konaç\n// This file is licensed under the MIT license.\n\n//go:build ignore\n\npackage main\n\nvar Color vec4\nvar ShadedCorners [4]float\nvar RectSize vec2\nvar RectPosInUnit vec2\nvar TotalSize vec2\n\nfunc Fragment(position vec4, texCoord vec2, color vec4) vec4 {\n\tnormColor := Color / 0xffff\n\n\tposInUnit := vec2(RectSize.x*texCoord.x, RectSize.y*texCoord.y) + RectPosInUnit\n\tif TotalSize.x <= TotalSize.y {\n\t\tradius := TotalSize.x / 2\n\t\troundCenter1 := vec2(radius)\n\t\troundCenter2 := vec2(radius, TotalSize.y-radius)\n\n\t\tdistToCenter1 := distance(posInUnit, roundCenter1)\n\t\tdistToCenter2 := distance(posInUnit, roundCenter2)\n\n\t\t// Top Left Corner\n\t\tif (ShadedCorners[0] > 0) && (posInUnit.x < radius) &&\n\t\t\t(posInUnit.y < radius) && (distToCenter1 > radius) {\n\t\t\tnormColor.a = 0\n\t\t}\n\t\t// Bottom Left Corner\n\t\tif (ShadedCorners[1] > 0) && (posInUnit.x < radius) &&\n\t\t\t(posInUnit.y > (TotalSize.y - radius)) && (distToCenter2 > radius) {\n\t\t\tnormColor.a = 0\n\t\t}\n\t\t// Bottom Right Corner\n\t\tif (ShadedCorners[2] > 0) && (posInUnit.x > radius) &&\n\t\t\t(posInUnit.y > (TotalSize.y - radius)) && (distToCenter2 > radius) {\n\t\t\tnormColor.a = 0\n\t\t}\n\t\t// Top Right Corner\n\t\tif (ShadedCorners[3] > 0) && (posInUnit.x > radius) &&\n\t\t\t(posInUnit.y < radius) && (distToCenter1 > radius) {\n\t\t\tnormColor.a = 0\n\t\t}\n\t} else {\n\t\tradius := TotalSize.y / 2\n\t\troundCenter1 := vec2(radius)\n\t\troundCenter2 := vec2(TotalSize.x-radius, radius)\n\n\t\tdistToCenter1 := distance(posInUnit, roundCenter1)\n\t\tdistToCenter2 := distance(posInUnit, roundCenter2)\n\n\t\t// Top Left Corner\n\t\tif (ShadedCorners[0] > 0) && (posInUnit.x < radius) &&\n\t\t\t(posInUnit.y < radius) && (distToCenter1 > radius) {\n\t\t\tnormColor.a = 0\n\t\t}\n\t\t// Bottom Left Corner\n\t\tif (ShadedCorners[1] > 0) && (posInUnit.x < radius) &&\n\t\t\t(posInUnit.y > radius) && (distToCenter1 > radius) {\n\t\t\tnormColor.a = 0\n\t\t}\n\t\t// Bottom Right Corner\n\t\tif (ShadedCorners[2] > 0) && (posInUnit.x > TotalSize.x-radius) &&\n\t\t\t(posInUnit.y > radius) && (distToCenter2 > radius) {\n\t\t\tnormColor.a = 0\n\t\t}\n\t\t// Top Right Corner\n\t\tif (ShadedCorners[3] > 0) && (posInUnit.x > TotalSize.x-radius) &&\n\t\t\t(posInUnit.y < radius) && (distToCenter2 > radius) {\n\t\t\tnormColor.a = 0\n\t\t}\n\t}\n\n\tnormColor.rgb *= normColor.a\n\treturn normColor\n}\n")
