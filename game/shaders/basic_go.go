// Code generated by file2byteslice. DO NOT EDIT.

package shaders

var Basic = []byte("//go:build ignore\r\n\r\npackage main\r\n\r\nfunc Fragment(position vec4, texCoord vec2, color vec4) vec4 {\r\n\tnormColor := color / 0xffff\r\n\tnormColor.rgb *= normColor.a\r\n\treturn normColor\r\n}\r\n")
