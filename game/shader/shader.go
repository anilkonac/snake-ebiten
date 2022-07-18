package shader

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	Basic     = "basic.kage.go"
	Round     = "round.kage.go"
	Score     = "score.kage.go"
	SnakeHead = "snakehead.kage.go"
	Title     = "title.kage.go"
)

var (
	//go:embed *.kage.go
	fs embed.FS
)

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
