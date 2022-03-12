package game

import (
	"github.com/anilkonac/snake-ebiten/game/shaders"
	"github.com/hajimehoshi/ebiten/v2"
)

type shaderT uint8

const (
	shaderBasic shaderT = iota
	shaderRound
	shaderHollow
	shaderTotal
)

var shaderMap map[shaderT]*ebiten.Shader
var curShader = shaderRound

func init() {
	shaderMap = make(map[shaderT]*ebiten.Shader)
	newShader(shaderBasic, shaders.Basic)
	newShader(shaderRound, shaders.Round)
	newShader(shaderHollow, shaders.Hollow)
}

func newShader(shaderNum shaderT, src []byte) {
	shader, err := ebiten.NewShader(src)
	if err != nil {
		panic(err)
	}
	shaderMap[shaderNum] = shader
}
