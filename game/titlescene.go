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

	"github.com/anilkonac/snake-ebiten/game/params"
	"github.com/anilkonac/snake-ebiten/game/shaders"
	s "github.com/anilkonac/snake-ebiten/game/snake"
	t "github.com/anilkonac/snake-ebiten/game/tools"
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
	dumbSnakeRunMultip  = 2.5
)

// Title Rectangle parameters
const (
	titleRectWidth                 = 540
	titleRectHeight                = 405
	titleRectRatio                 = 1.0 * titleRectWidth / titleRectHeight
	titleRectCornerRadiusX         = params.RadiusSnake
	titleRectCornerRadiusY         = titleRectCornerRadiusX / titleRectRatio
	titleRectInitialAlpha          = 230 / 255.0
	titleRectDissapearRate float32 = (85 / 255.0) * params.DeltaTime
	textTitle                      = "Ssnake"
	textPressToPlay                = "Press any key to start"
	textTitleShiftY                = -50
	textKeyPromptShiftY            = +100
	keyPromptShowTime              = 1.0 //sec
	keyPromptHideTime              = 0.5 // sec
)

var (
	titleSceenAlive        = true
	colorTitleRect         = &params.ColorSnake2
	titleBackgroundIndices = []uint16{
		1, 0, 2,
		2, 3, 1,
	}

	snakeColors = map[int]*color.RGBA{
		0: &params.ColorSnake1,
		1: &params.ColorSnake2,
		2: &params.ColorFood,
		3: &params.ColorDebug,
	}
)

type titleScene struct {
	titleRectAlpha      float32
	snakes              []*s.Snake
	pressedKeys         []ebiten.Key
	titleRectVertices   []ebiten.Vertex
	titleImage          *ebiten.Image
	titleImageKeyPrompt *ebiten.Image
	titleRectDrawOpts   ebiten.DrawTrianglesShaderOptions
	shaderTitle         *ebiten.Shader
}

func newTitleScene() *titleScene {
	// Create title rect model
	titleRect := t.RectF32{
		Pos:       t.Vec32{X: (params.ScreenWidth - titleRectWidth) / 2.0, Y: (params.ScreenHeight - titleRectHeight) / 2.0},
		Size:      t.Vec32{X: titleRectWidth, Y: titleRectHeight},
		PosInUnit: t.Vec32{X: 0, Y: 0},
	}

	// Create scene
	scene := &titleScene{
		titleRectAlpha:    titleRectInitialAlpha,
		snakes:            make([]*s.Snake, maxSnakes),
		pressedKeys:       make([]ebiten.Key, 0, 10),
		shaderTitle:       t.NewShader(shaders.Title),
		titleRectVertices: titleRect.Vertices(colorTitleRect),
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
		snake := s.NewSnakeRandDirLoc(uint16(length), speed, snakeColors[rand.Intn(lenSnakeColors)])
		scene.snakes[iSnake] = snake

		// Activate snake bot
		go controlDumbly(snake)
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
		params.ColorBackground)

	// Prepare key prompt text image
	t.titleImageKeyPrompt = ebiten.NewImageFromImage(t.titleImage)

	// Draw key prompt text to the image
	text.Draw(t.titleImageKeyPrompt, textPressToPlay, fontFaceScore,
		(titleRectWidth-boundTextKeyPromptSize.X)/2.0-boundTextKeyPrompt.Min.X,
		(titleRectHeight-boundTextKeyPromptSize.Y)/2.0-boundTextKeyPrompt.Min.Y+textKeyPromptShiftY, params.ColorBackground)

	// Send images to the shader
	t.titleRectDrawOpts.Images[0] = t.titleImage
	t.titleRectDrawOpts.Images[1] = t.titleImageKeyPrompt

	go t.keyPromptFlipFlop()
}

func (t *titleScene) update() bool {
	for _, snake := range t.snakes {
		snake.Update(params.EatingAnimStartDistance) // Make sure the snake's mouth is not open
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
		params.TeleportActive = false
		t.titleRectDrawOpts.Uniforms["ShowKeyPrompt"] = float32(0.0)

		for _, snake := range t.snakes {
			snake.Speed *= dumbSnakeRunMultip
		}
	}
}

func (t *titleScene) draw(screen *ebiten.Image) {
	screen.Fill(params.ColorBackground)

	// Draw snakes
	for _, snake := range t.snakes {
		for unit := snake.UnitHead; unit != nil; unit = unit.Next {
			draw(screen, unit)
		}
	}

	drawFPS(screen)

	// Draw Title Rect
	screen.DrawTrianglesShader(t.titleRectVertices, titleBackgroundIndices, t.shaderTitle, &t.titleRectDrawOpts)

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
func controlDumbly(snake *s.Snake) {
	const turnTimeMinMs = turnTimeMin * 1000
	const turnTimeDiffMs = turnTimeDiff * 1000

	var dirNew s.DirectionT
	for titleSceenAlive {
		// Determine the new direction.
		dirCurrent := snake.LastDirection()

		randResult := rand.Float32()
		if dirCurrent.IsVertical() {
			if randResult < 0.5 {
				dirNew = s.DirectionLeft
			} else {
				dirNew = s.DirectionRight
			}
		} else {
			if randResult < 0.5 {
				dirNew = s.DirectionUp
			} else {
				dirNew = s.DirectionDown
			}
		}

		// Create a new turn and take it
		newTurn := s.NewTurn(dirCurrent, dirNew)
		snake.TurnTo(newTurn, false)

		// Sleep for a random time limited by turnTimeMax and turnTimeMin.
		sleepTime := time.Duration(turnTimeMinMs + rand.Float32()*turnTimeDiffMs)
		time.Sleep(time.Millisecond * sleepTime)
	}
}
