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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	titleRectWidth                 = 540
	titleRectHeight                = 405
	titleRectRatio                 = 1.0 * titleRectWidth / titleRectHeight
	titleRectCornerRadiusX         = radiusSnake
	titleRectCornerRadiusY         = titleRectCornerRadiusX / titleRectRatio
	titleRectInitialAlpha          = 230 / 255.0
	titleRectDissapearRate float32 = (75 / 255.0) * deltaTime
	textTitle                      = "Ssnake"
	textPressToPlay                = "Press any key to start"
	textTitleShiftY                = -50
	textKeyPromptShiftY            = +100
	keyPromptShowTime              = 1.0 //sec
	keyPromptHideTime              = 0.5 // sec
)

var (
	titleSceenAlive        = true
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
)

type titleScene struct {
	titleRectAlpha      float32
	snakes              []*snake
	pressedKeys         []ebiten.Key
	titleRectVertices   []ebiten.Vertex
	titleImage          *ebiten.Image
	titleImageKeyPrompt *ebiten.Image
	titleRectDrawOpts   ebiten.DrawTrianglesShaderOptions
}

func newTitleScreen() *titleScene {
	// Create title rect model
	titleRect := rectF32{
		pos:       vec32{(ScreenWidth - titleRectWidth) / 2.0, (ScreenHeight - titleRectHeight) / 2.0},
		size:      vec32{titleRectWidth, titleRectHeight},
		posInUnit: vec32{0, 0},
	}

	// Create scene
	scene := &titleScene{
		snakes:            make([]*snake, maxSnakes),
		titleRectVertices: titleRect.vertices(colorTitleRect),
		titleRectAlpha:    titleRectInitialAlpha,
		pressedKeys:       make([]ebiten.Key, 0, 10),
		titleRectDrawOpts: ebiten.DrawTrianglesShaderOptions{
			Uniforms: map[string]interface{}{
				"ShowKeyPrompt": float32(0.0),
				"RadiusTex":     []float32{float32(titleRectCornerRadiusX / titleRectWidth), float32(titleRectCornerRadiusY / titleRectHeight)},
				"Alpha":         float32(titleRectInitialAlpha),
			},
		},
	}
	scene.prepareTitleRects()

	// Create snakes
	lenSnakeColors := len(snakeColors)
	for iSnake := 0; iSnake < maxSnakes; iSnake++ {
		length := dumbSnakeLengthMin + rand.Intn(dumbSnakeLengthDiff)
		speed := dumbSnakeSpeedMin + rand.Float64()*dumbSnakeSpeedDiff
		snake := newSnakeRandDirLoc(uint16(length), speed, snakeColors[rand.Intn(lenSnakeColors)])
		scene.snakes[iSnake] = snake

		// Activate snake bot
		go snake.controlDumbly()
	}

	return scene
}

func (t *titleScene) prepareTitleRects() {
	// Prepare title text image
	t.titleImage = ebiten.NewImage(titleRectWidth, titleRectHeight)
	t.titleImage.Fill(colorTitleRect)

	boundTextTitleSize := boundTextTitle.Size()
	boundTextKeyPromptSize := boundTextKeyPrompt.Size()

	// Draw Title text to the image
	text.Draw(t.titleImage, textTitle, fontFaceTitle,
		(titleRectWidth-boundTextTitleSize.X)/2.0-boundTextTitle.Min.X,
		(titleRectHeight-boundTextTitleSize.Y)/2.0-boundTextTitle.Min.Y+textTitleShiftY,
		colorBackground)

	// Prepare key prompt text image
	t.titleImageKeyPrompt = ebiten.NewImageFromImage(t.titleImage)

	// Draw key prompt text to the image
	text.Draw(t.titleImageKeyPrompt, textPressToPlay, fontFaceScore,
		(titleRectWidth-boundTextKeyPromptSize.X)/2.0-boundTextKeyPrompt.Min.X,
		(titleRectHeight-boundTextKeyPromptSize.Y)/2.0-boundTextKeyPrompt.Min.Y+textKeyPromptShiftY, colorBackground)

	// Send the images to the shader
	t.titleRectDrawOpts.Images[0] = t.titleImage
	t.titleRectDrawOpts.Images[1] = t.titleImageKeyPrompt

	go t.keyPromptFlipFlop()
}

func (t *titleScene) update() bool {
	for _, snake := range t.snakes {
		snake.update(eatingAnimStartDistance) // Make sure the snake's mouth is not open
	}

	t.handleKeyPress()

	// Update transition process to the next scene
	if !titleSceenAlive {
		t.titleRectAlpha -= titleRectDissapearRate
		t.titleRectDrawOpts.Uniforms["Alpha"] = t.titleRectAlpha
		if t.titleRectAlpha <= 0.0 {
			return true
		}
	}

	return false
}

func (t *titleScene) handleKeyPress() {
	t.pressedKeys = inpututil.AppendPressedKeys(t.pressedKeys[:0])
	if len(t.pressedKeys) > 0 && titleSceenAlive {
		// Start transition process
		titleSceenAlive = false
		teleportActive = false
		t.titleRectDrawOpts.Uniforms["ShowKeyPrompt"] = float32(0.0)

		for _, snake := range t.snakes {
			snake.speed *= 2
		}
	}
}

func (t *titleScene) draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)

	// Draw snakes
	for _, snake := range t.snakes {
		for unit := snake.unitHead; unit != nil; unit = unit.next {
			draw(screen, unit)
		}
	}

	drawFPS(screen)

	// Draw Title Rect
	screen.DrawTrianglesShader(t.titleRectVertices, titleBackgroundIndices, shaderTitle, &t.titleRectDrawOpts)

}

// Goroutine
func (t *titleScene) keyPromptFlipFlop() {
	const showTimeMs = keyPromptShowTime * 1000
	const hideTimeMs = keyPromptHideTime * 1000

	showPrompt := true
	for titleSceenAlive {
		if showPrompt {
			t.titleRectDrawOpts.Uniforms["ShowKeyPrompt"] = float32(1.0)
			time.Sleep(time.Millisecond * time.Duration(showTimeMs))
		} else {
			t.titleRectDrawOpts.Uniforms["ShowKeyPrompt"] = float32(0.0)
			time.Sleep(time.Millisecond * time.Duration(hideTimeMs))
		}
		showPrompt = !showPrompt
	}
}

// Goroutine
// Dumb snake bot
func (s *snake) controlDumbly() {
	const turnTimeMinMs = turnTimeMin * 1000
	const turnTimeDiffMs = turnTimeDiff * 1000

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
		sleepTime := time.Duration(turnTimeMinMs + rand.Float32()*turnTimeDiffMs)
		time.Sleep(time.Millisecond * sleepTime)
	}
}
