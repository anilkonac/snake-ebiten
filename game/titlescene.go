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

// Dumb snake parameters
const (
	maxSnakes           = 30
	turnTimeMin         = 0.0 // sec
	turnTimeMax         = 1.5 // sec
	turnTimeDiff        = turnTimeMax - turnTimeMin
	dumbSnakeLengthMin  = 120
	dumbSnakeLengthMax  = 480
	dumbSnakeLengthDiff = dumbSnakeLengthMax - dumbSnakeLengthMin
	dumbSnakeSpeedMin   = 150
	dumbSnakeSpeedMax   = 400
	dumbSnakeSpeedDiff  = dumbSnakeSpeedMax - dumbSnakeSpeedMin
)

// Title Rectangle parameters
const (
	titleRectWidth         = 540
	titleRectHeight        = 405
	titleRectRatio         = 1.0 * titleRectWidth / titleRectHeight
	titleRectCornerRadiusX = radiusSnake
	titleRectCornerRadiusY = titleRectCornerRadiusX / titleRectRatio
	titleRectInitialAlpha  = 230 / 255.0
	textTitle              = "Ssnake"
	textPressToPlay        = "Press any key to start"
	textTitleShiftY        = -50
	textKeyPromptShiftY    = +100
	keyPromptShowTime      = 1.0 //sec
	keyPromptHideTime      = 0.5 // sec
)

const (
	titleRectDissapearRate float32 = (75 / 255.0) * deltaTime
)

var (
	titleSceenAlive        = true
	titleImage             *ebiten.Image
	titleImageKeyPrompt    *ebiten.Image
	colorTitleRect         = &colorSnake2
	titleBackgroundIndices = []uint16{
		1, 0, 2,
		2, 3, 1,
	}

	snakeColors = map[int]*color.RGBA{
		0: &colorSnake1,
		1: &colorSnake2,
		2: &colorFood,
		3: &colorDebug,
	}

	titleShaderOp = ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]interface{}{
			"ShowKeyPrompt": float32(0.0),
			"RadiusTex":     []float32{float32(titleRectCornerRadiusX / titleRectWidth), float32(titleRectCornerRadiusY / titleRectHeight)},
			"Alpha":         float32(titleRectInitialAlpha),
		},
	}
)

type titleScene struct {
	snakes                  []*snake
	titleBackground         rectF32
	titleBackgroundVertices []ebiten.Vertex
	titleRectAlpha          float32
}

func newTitleScreen() *titleScene {
	scene := &titleScene{
		snakes:          make([]*snake, maxSnakes),
		titleBackground: *newRect(vec32{(ScreenWidth - titleRectWidth) / 2.0, (ScreenHeight - titleRectHeight) / 2.0}, vec32{titleRectWidth, titleRectHeight}),
		titleRectAlpha:  titleRectInitialAlpha,
	}
	scene.titleBackgroundVertices = scene.titleBackground.vertices(colorTitleRect)

	lenSnakeColors := len(snakeColors)
	for iSnake := 0; iSnake < maxSnakes; iSnake++ {
		length := dumbSnakeLengthMin + rand.Intn(dumbSnakeLengthDiff)
		speed := dumbSnakeSpeedMin + rand.Float64()*dumbSnakeSpeedDiff
		snake := newSnakeRandDirLoc(uint16(length), speed, snakeColors[rand.Intn(lenSnakeColors)])
		go snake.controlDumbly()
		scene.snakes[iSnake] = snake
	}
	return scene
}

func initTitleRect() {
	// Prepare title text image
	titleImage = ebiten.NewImage(titleRectWidth, titleRectHeight)
	titleImage.Fill(colorTitleRect)

	boundTextTitleSize := boundTextTitle.Size()
	boundTextKeyPromptSize := boundTextKeyPrompt.Size()
	// Draw Title
	text.Draw(titleImage, textTitle, fontFaceTitle,
		(titleRectWidth-boundTextTitleSize.X)/2.0-boundTextTitle.Min.X,
		(titleRectHeight-boundTextTitleSize.Y)/2.0-boundTextTitle.Min.Y+textTitleShiftY,
		colorBackground)

	titleImageKeyPrompt = ebiten.NewImageFromImage(titleImage)

	// Draw key prompt
	text.Draw(titleImageKeyPrompt, textPressToPlay, fontFaceScore,
		(titleRectWidth-boundTextKeyPromptSize.X)/2.0-boundTextKeyPrompt.Min.X,
		(titleRectHeight-boundTextKeyPromptSize.Y)/2.0-boundTextKeyPrompt.Min.Y+textKeyPromptShiftY, colorBackground)

	titleShaderOp.Images[0] = titleImage
	titleShaderOp.Images[1] = titleImageKeyPrompt

	go keyPromptFlipFlop()
}

func (t *titleScene) update(input *inputHandler) bool {
	for _, snake := range t.snakes {
		snake.update(eatingAnimStartDistance)
	}

	if len(input.keys) > 0 && titleSceenAlive {
		titleSceenAlive = false
		teleportActive = false
		for _, snake := range t.snakes {
			snake.speed *= 2
		}
	}

	if !titleSceenAlive {
		t.titleRectAlpha -= titleRectDissapearRate
		titleShaderOp.Uniforms["Alpha"] = t.titleRectAlpha
		if t.titleRectAlpha <= 0.0 {
			return true
		}
		titleShaderOp.Uniforms["ShowKeyPrompt"] = float32(0.0)
	}
	return false
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
