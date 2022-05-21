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
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	maxSnakes                     = 30
	turnTimeMin                   = 0.0 // sec
	turnTimeMax                   = 2.0 // sec
	turnTimeDiff                  = turnTimeMax - turnTimeMin
	titleBackgroundWidth  float32 = 600
	titleBackgroundHeight float32 = 400
)

var (
	dumbSnakesAlive = true
	snakeColors     = map[int]*color.RGBA{
		0: &colorSnake1,
		1: &colorSnake2,
		2: &colorFood,
		3: &colorDebug,
	}
	colorTitleScreen  = color.RGBA{colorSnake1.R, colorSnake1.G, colorSnake1.B, 230}
	titleBackgroundOp = ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]interface{}{
			"Radius":    float32(titleBackgroundHeight / 2.0),
			"Size":      []float32{titleBackgroundWidth, titleBackgroundHeight},
			"Direction": float32(2),
		},
	}
	titleBackgroundIndices = []uint16{
		1, 0, 2,
		2, 3, 1,
	}
)

type titleScene struct {
	snakes                  []*snake
	titleBackground         rectF32
	titleBackgroundVertices []ebiten.Vertex
}

func newTitleScreen() *titleScene {
	scene := &titleScene{
		snakes:          make([]*snake, maxSnakes),
		titleBackground: *newRect(vec32{(ScreenWidth - titleBackgroundWidth) / 2.0, (ScreenHeight - titleBackgroundHeight) / 2.0}, vec32{titleBackgroundWidth, titleBackgroundHeight}),
	}
	scene.titleBackgroundVertices = scene.titleBackground.vertices(&colorTitleScreen)

	lenSnakeColors := len(snakeColors)
	for iSnake := 0; iSnake < maxSnakes; iSnake++ {
		snake := newSnakeRandDirLoc(snakeColors[rand.Intn(lenSnakeColors)])
		go snake.controlDumbly()
		scene.snakes[iSnake] = snake
	}
	return scene
}

func (t *titleScene) update() {
	for _, snake := range t.snakes {
		snake.update(eatingAnimStartDistance)
	}
}

func (t *titleScene) draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)

	// Draw snakes
	for _, snake := range t.snakes {
		// Draw the snake
		for unit := snake.unitHead; unit != nil; unit = unit.next {
			draw(screen, unit)
		}
	}

	// Draw Title Background
	screen.DrawTrianglesShader(t.titleBackgroundVertices, titleBackgroundIndices, shaderRound, &titleBackgroundOp)

	// Draw Title
	text.Draw(screen, "Ssnake", fontFaceScore, halfScreenWidth, halfScreenHeight, colorBackground)

	// Draw key prompt
	text.Draw(screen, "Press any key to play", fontFaceScore, halfScreenWidth-100, halfScreenHeight+100, colorBackground)
}

// Goroutine
func (s *snake) controlDumbly() {
	var dirNew directionT
	for dumbSnakesAlive {
		// Determine the new direction.
		dirCurrent := s.lastDirection()

		randResult := rand.Float32()
		if dirCurrent.isVertical() {
			if randResult < 0.5 {
				dirNew = directionLeft
			} else {
				dirNew = directionRight
			}
		} else {
			if randResult < 0.5 {
				dirNew = directionUp
			} else {
				dirNew = directionDown
			}
		}

		// Create a new turn and take it
		newTurn := newTurn(dirCurrent, dirNew)
		s.turnTo(newTurn, false)

		// Sleep for a random time limited by turnTimeMax and turnTimeMin.
		sleepTime := time.Duration((turnTimeMin + rand.Float32()*turnTimeDiff) * 1000)
		time.Sleep(time.Millisecond * sleepTime)
	}
}
