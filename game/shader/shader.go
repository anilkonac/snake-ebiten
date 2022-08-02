package shader

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	PathBasic       = "basic.kage.go"
	PathCircle      = "circle.kage.go"
	PathCircleMouth = "circlemouth.kage.go"
	PathTitle       = "title.kage.go"
)

var (
	//go:embed *.kage.go
	fs     embed.FS
	Circle ebiten.Shader
)

func init() {
	Circle = *New(PathCircle)
}

func New(sh string) *ebiten.Shader {
	bytesShader, err := fs.ReadFile(sh)
	if err != nil {
		panic(err)
	}

	shader, err := ebiten.NewShader(bytesShader)
	if err != nil {
		panic(err)
	}

	return shader
}
