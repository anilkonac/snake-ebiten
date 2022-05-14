package game

import (
	"github.com/anilkonac/snake-ebiten/game/shaders"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	shaderRound     *ebiten.Shader
	shaderScore     *ebiten.Shader
	shaderSnakeHead *ebiten.Shader
)

func init() {
	shaderRound = newShader(shaders.Round)
	shaderScore = newShader(shaders.Score)
	shaderSnakeHead = newShader(shaders.SnakeHead)
}

func newShader(src []byte) *ebiten.Shader {
	shader, err := ebiten.NewShader(src)
	if err != nil {
		panic(err)
	}
	return shader
}
