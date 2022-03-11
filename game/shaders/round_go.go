// Code generated by file2byteslice. DO NOT EDIT.

package shaders

var Round = []byte("// Copyright (c) 2022 Anıl Konaç\n// This file is licensed under the MIT license.\n\n//go:build ignore\n\npackage main\n\nvar Color vec4\nvar Width float\nvar Height float\n\nfunc Fragment(position vec4, texCoord vec2, color vec4) vec4 {\n\tnormColor := Color / 0xffff\n\n\trealCoord := vec2(Width*texCoord.x, Height*texCoord.y)\n\tif Width <= Height {\n\t\tradius := Width / 2\n\t\troundCenter1 := vec2(radius, radius)\n\t\troundCenter2 := vec2(radius, Height-radius)\n\n\t\tdistToCenter1 := distance(realCoord, roundCenter1)\n\t\tdistToCenter2 := distance(realCoord, roundCenter2)\n\t\tif ((realCoord.y < radius) && (distToCenter1 > radius)) ||\n\t\t\t((realCoord.y > (Height - radius)) && (distToCenter2 > radius)) {\n\t\t\tnormColor.a = 0\n\t\t}\n\t} else {\n\t\tradius := Height / 2\n\t\troundCenter1 := vec2(radius)\n\t\troundCenter2 := vec2(Width-radius, radius)\n\n\t\tdistToCenter1 := distance(realCoord, roundCenter1)\n\t\tdistToCenter2 := distance(realCoord, roundCenter2)\n\t\tif ((realCoord.x < radius) && (distToCenter1 > radius)) ||\n\t\t\t((realCoord.x > (Width - radius)) && (distToCenter2 > radius)) {\n\t\t\tnormColor.a = 0\n\t\t}\n\t}\n\n\tnormColor.xyz *= normColor.a\n\treturn normColor\n}\n")
