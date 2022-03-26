//go:build ignore

package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	normColor := color / 0xffff
	normColor.rgb *= normColor.a
	return normColor
}
