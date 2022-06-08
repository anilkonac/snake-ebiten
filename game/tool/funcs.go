package tool

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

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

func MarkPoint(dst *ebiten.Image, p Vec64, length float64, clr color.Color) {
	ebitenutil.DrawLine(dst, p.X-length, p.Y, p.X+length, p.Y, clr)
	ebitenutil.DrawLine(dst, p.X, p.Y-length, p.X, p.Y+length, clr)
}
