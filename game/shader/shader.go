package shader

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
)

type shaderE uint8

const (
	Basic shaderE = iota
	Round
	Score
	SnakeHead
	Title
)

var (
	//go:embed *.kage
	fs            embed.FS
	mapShaderPath = map[shaderE]string{
		Basic:     "basic.kage",
		Round:     "round.kage",
		Score:     "score.kage",
		SnakeHead: "snakehead.kage",
		Title:     "title.kage",
	}
)

func New(sh shaderE) *ebiten.Shader {
	bytesShader, err := fs.ReadFile(mapShaderPath[sh])
	if err != nil {
		panic(err)
	}

	shader, err := ebiten.NewShader(bytesShader)
	if err != nil {
		panic(err)
	}

	return shader
}
