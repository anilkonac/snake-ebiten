package game

import (
	"github.com/anilkonac/snake-ebiten/game/params"
	"github.com/anilkonac/snake-ebiten/game/shaders"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	shaderScore     *ebiten.Shader
	shaderSnakeHead *ebiten.Shader
	shaderTitle     *ebiten.Shader
)

func init() {
	params.ShaderRound = newShader(shaders.Round)
	shaderScore = newShader(shaders.Score)
	shaderSnakeHead = newShader(shaders.SnakeHead)
	shaderTitle = newShader(shaders.Title)
}

func newShader(src []byte) *ebiten.Shader {
	shader, err := ebiten.NewShader(src)
	if err != nil {
		panic(err)
	}
	return shader
}
