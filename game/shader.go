/*
snake-ebiten
Copyright (C) 2022 Anıl Konaç

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package game

import (
	"github.com/anilkonac/snake-ebiten/game/shaders"
	"github.com/hajimehoshi/ebiten/v2"
)

type shaderT uint8

const (
	shaderBasic shaderT = iota
	// shaderHollow
	shaderRound
	shaderTotal
)

var shaderMap map[shaderT]*ebiten.Shader

var curShader = shaderRound

func init() {
	shaderMap = make(map[shaderT]*ebiten.Shader)
	newShader(shaderBasic, shaders.Basic)
	newShader(shaderRound, shaders.Round)
	// newShader(shaderHollow, shaders.Hollow)
}

func newShader(shaderNum shaderT, src []byte) {
	shader, err := ebiten.NewShader(src)
	if err != nil {
		panic(err)
	}
	shaderMap[shaderNum] = shader
}
