// Code generated by file2byteslice. DO NOT EDIT.

package shaders

var Basic = []byte("//go:build ignore\n\npackage main\n\nvar Color vec4\n\nfunc Fragment(position vec4, texCoord vec2, color vec4) vec4 {\n\tnormColor := Color / 0xffff\n\tnormColor.rgb *= normColor.a\n\treturn normColor\n}\n")
