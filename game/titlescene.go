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
	maxSnakes           = 30
	turnTimeMin         = 0.0 // sec
	turnTimeMax         = 1.5 // sec
	turnTimeDiff        = turnTimeMax - turnTimeMin
	snakeLengthMin      = 120
	snakeLengthMax      = 480
	snakeLengthDiff     = snakeLengthMax - snakeLengthMin
	titleRectWidth      = 600
	titleRectHeight     = 400
	textTitle           = "Ssnake"
	textPressToPlay     = "Press any key to start"
	textTitleShiftX     = 0
	textTitleShiftY     = -50
	textKeyPromptShiftX = 0
	textKeyPromptShiftY = +100
	keyPromptShowTime   = 1.0 //sec
	keyPromptHideTime   = 0.5 // sec
)

var (
	titleSceenAlive = true
	snakeColors     = map[int]*color.RGBA{
		0: &colorSnake1,
		1: &colorSnake2,
		2: &colorFood,
		3: &colorDebug,
	}
	colorTitleScreen = color.RGBA{colorSnake1.R, colorSnake1.G, colorSnake1.B, 230}
	titleShaderOp    = ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]interface{}{
			"Radius":        float32(titleRectHeight / 2.0),
			"Size":          []float32{float32(titleRectWidth), float32(titleRectHeight)},
			"Direction":     float32(2),
			"ShowKeyPrompt": float32(0.0),
		},
	}
	titleBackgroundIndices = []uint16{
		1, 0, 2,
		2, 3, 1,
	}
	titleImage          *ebiten.Image
	titleImageKeyPrompt *ebiten.Image
)

type titleScene struct {
	snakes                  []*snake
	titleBackground         rectF32
	titleBackgroundVertices []ebiten.Vertex
}

func newTitleScreen() *titleScene {
	scene := &titleScene{
		snakes:          make([]*snake, maxSnakes),
		titleBackground: *newRect(vec32{(ScreenWidth - titleRectWidth) / 2.0, (ScreenHeight - titleRectHeight) / 2.0}, vec32{titleRectWidth, titleRectHeight}),
	}
	scene.titleBackgroundVertices = scene.titleBackground.vertices(&colorTitleScreen)

	lenSnakeColors := len(snakeColors)
	for iSnake := 0; iSnake < maxSnakes; iSnake++ {
		length := snakeLengthMin + rand.Intn(snakeLengthDiff)
		snake := newSnakeRandDirLoc(uint16(length), snakeColors[rand.Intn(lenSnakeColors)])
		go snake.controlDumbly()
		scene.snakes[iSnake] = snake
	}
	return scene
}

func initTitle() {
	// Prepare title text image
	titleImage = ebiten.NewImage(titleRectWidth, titleRectHeight)
	// titleImageKeyPrompt = ebiten.NewImage(titleRectWidth, titleRectHeight)
	titleImage.Fill(colorSnake2)

	println(boundTextTitle.Min.X)
	boundTextTitleSize := boundTextTitle.Size()
	boundTextKeyPromptSize := boundTextKeyPrompt.Size()
	// Draw Title
	text.Draw(titleImage, textTitle, fontFaceTitle,
		(titleRectWidth-boundTextTitleSize.X)/2.0-boundTextTitle.Min.X+textTitleShiftX,
		(titleRectHeight-boundTextTitleSize.Y)/2.0-boundTextTitle.Min.Y+textTitleShiftY,
		colorBackground)

	titleImageKeyPrompt = ebiten.NewImageFromImage(titleImage)

	// Draw key prompt
	text.Draw(titleImageKeyPrompt, textPressToPlay, fontFaceScore,
		(titleRectWidth-boundTextKeyPromptSize.X)/2.0-boundTextKeyPrompt.Min.X+textKeyPromptShiftX,
		(titleRectHeight-boundTextKeyPromptSize.Y)/2.0-boundTextKeyPrompt.Min.Y+textKeyPromptShiftY, colorBackground)

	titleShaderOp.Images[0] = titleImage
	titleShaderOp.Images[1] = titleImageKeyPrompt

	go keyPromptFlipFlop()
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
	screen.DrawTrianglesShader(t.titleBackgroundVertices, titleBackgroundIndices, shaderTitle, &titleShaderOp)

}

func (t *titleScene) exit() {
	titleSceenAlive = false
}

// Goroutine
func keyPromptFlipFlop() {
	showPrompt := true
	for titleSceenAlive {
		if showPrompt {
			titleShaderOp.Uniforms["ShowKeyPrompt"] = float32(1.0)
			time.Sleep(time.Millisecond * time.Duration(keyPromptShowTime*1000))
		} else {
			titleShaderOp.Uniforms["ShowKeyPrompt"] = float32(0.0)
			time.Sleep(time.Millisecond * time.Duration(keyPromptHideTime*1000))
		}
		showPrompt = !showPrompt
	}
}

// Goroutine
func (s *snake) controlDumbly() {
	var dirNew directionT
	for titleSceenAlive {
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
