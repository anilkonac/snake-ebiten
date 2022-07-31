/*
Copyright (C) 2022 Anıl Konaç

This file is part of snake-ebiten.

snake-ebiten is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

snake-ebiten is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with snake-ebiten. If not, see <https://www.gnu.org/licenses/>.
*/

package game

import (
	"github.com/anilkonac/snake-ebiten/game/object/snake"
	"github.com/anilkonac/snake-ebiten/game/param"
	"github.com/hajimehoshi/ebiten/v2"
)

type scene interface {
	update() bool // Return true if the scene is finished
	draw(*ebiten.Image)
}

// Game implements ebiten.Game interface.
type Game struct {
	curScene    scene
	playerSnake *snake.Snake
}

func NewGame() *Game {
	playerSnake := snake.NewSnakeRandDirLoc(param.SnakeLength, param.SnakeSpeedInitial, &param.ColorSnake1)

	return &Game{
		curScene: newTitleScene(playerSnake),
		// curScene:    newGameScene(playerSnake),
		playerSnake: playerSnake,
	}
}

// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	if g.curScene.update() {
		switch g.curScene.(type) {
		case *titleScene:
			g.curScene = newGameScene(g.playerSnake)
		}
	}

	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	g.curScene.draw(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return param.ScreenWidth, param.ScreenHeight
}

func ScreenSize() (int, int) {
	return param.ScreenWidth, param.ScreenHeight
}
