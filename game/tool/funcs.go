package tool

import "github.com/hajimehoshi/ebiten/v2"

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func NewShader(src []byte) *ebiten.Shader {
	shader, err := ebiten.NewShader(src)
	Panic(err)
	return shader
}