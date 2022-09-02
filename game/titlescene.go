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
	"image/color"
	"math/rand"
	"time"

	c "github.com/anilkonac/snake-ebiten/game/core"
	s "github.com/anilkonac/snake-ebiten/game/object/snake"
	"github.com/anilkonac/snake-ebiten/game/param"
	"github.com/anilkonac/snake-ebiten/game/shader"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// Dumb snake parameters
const (
	numBotSnakes        = 29
	turnTimeMinSec      = 0
	turnTimeMaxSec      = 2
	turnTimeDiff        = turnTimeMaxSec - turnTimeMinSec
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
	titleRectCornerRadiusX         = param.RadiusSnake
	titleRectCornerRadiusY         = titleRectCornerRadiusX / titleRectRatio
	titleRectInitialAlpha          = 230 / 255.0
	titleRectDissapearRate float32 = (80 / 255.0) * param.DeltaTime
	textTitle                      = "Ssnake"
	textPressToPlay                = "Press any key to start"
	textTitleShiftY                = -50
	textKeyPromptShiftY            = +100
	keyPromptShowTimeSec           = 1.0
	keyPromptHideTimeSec           = 0.5
)

var (
	titleSceneAlive = true
	colorTitleRect  = &param.ColorSnake2

	snakeColors = [...]*color.RGBA{
		0: &param.ColorSnake1,
		1: &param.ColorSnake2,
		2: &param.ColorFood,
		3: &param.ColorDebug,
	}
)

type titleScene struct {
	titleRectComp     c.TeleCompTriang
	titleRectAlpha    float32
	playerSnake       *s.Snake
	snakes            []s.Snake
	pressedKeys       []ebiten.Key
	shaderTitle       *ebiten.Shader
	titleRectDrawOpts ebiten.DrawTrianglesShaderOptions
}

func newTitleScene(playerSnake *s.Snake) *titleScene {
	// Create title rect model
	titleRect := c.RectF32{
		Pos:       c.Vec32{X: (param.ScreenWidth - titleRectWidth) / 2.0, Y: (param.ScreenHeight - titleRectHeight) / 2.0},
		Size:      c.Vec32{X: titleRectWidth, Y: titleRectHeight},
		PosInUnit: c.Vec32{X: 0, Y: 0},
	}

	// Create scene
	scene := &titleScene{
		playerSnake:    playerSnake,
		titleRectAlpha: titleRectInitialAlpha,
		snakes:         make([]s.Snake, 0, numBotSnakes),
		pressedKeys:    make([]ebiten.Key, 0, 10),
		shaderTitle:    shader.New(shader.PathTitle),
		titleRectDrawOpts: ebiten.DrawTrianglesShaderOptions{
			Uniforms: map[string]interface{}{
				"ShowKeyPrompt": float32(0.0),
				"RadiusTex":     []float32{float32(titleRectCornerRadiusX / titleRectWidth), float32(titleRectCornerRadiusY / titleRectHeight)},
				"Alpha":         float32(titleRectInitialAlpha),
			},
		},
	}
	scene.titleRectComp.SetColor(colorTitleRect)
	scene.titleRectComp.Update(&titleRect)
	scene.prepareTitleRects()

	// Create snakes
	// -------------

	// Create temp snakes for the title screen
	lenSnakeColors := len(snakeColors)
	for iSnake := 0; iSnake < numBotSnakes; iSnake++ {
		length := dumbSnakeLengthMin + rand.Intn(dumbSnakeLengthDiff)
		speed := dumbSnakeSpeedMin + rand.Float64()*dumbSnakeSpeedDiff
		scene.snakes = append(scene.snakes, *s.NewSnakeRandDirLoc(uint16(length), speed, snakeColors[rand.Intn(lenSnakeColors)]))

		go control(&scene.snakes[iSnake])

	}

	go control(playerSnake)

	return scene
}

func (t *titleScene) prepareTitleRects() {
	boundTextTitleSize := boundTextTitle.Size()
	boundTextKeyPromptSize := boundTextKeyPrompt.Size()

	// Prepare title text image
	titleImage := ebiten.NewImage(titleRectWidth, titleRectHeight)
	titleImage.Fill(colorTitleRect)

	// Draw Title text to the image
	text.Draw(titleImage, textTitle, fontFaceTitle,
		(titleRectWidth-boundTextTitleSize.X)/2.0-boundTextTitle.Min.X,
		(titleRectHeight-boundTextTitleSize.Y)/2.0-boundTextTitle.Min.Y+textTitleShiftY,
		param.ColorBackground)

	// Prepare key prompt text image
	titleImageKeyPrompt := ebiten.NewImageFromImage(titleImage)

	// Draw key prompt text to the image
	text.Draw(titleImageKeyPrompt, textPressToPlay, param.FontFaceScore,
		(titleRectWidth-boundTextKeyPromptSize.X)/2.0-boundTextKeyPrompt.Min.X,
		(titleRectHeight-boundTextKeyPromptSize.Y)/2.0-boundTextKeyPrompt.Min.Y+textKeyPromptShiftY, param.ColorBackground)

	// Send images to the shader
	t.titleRectDrawOpts.Images[0] = titleImage
	t.titleRectDrawOpts.Images[1] = titleImageKeyPrompt

	go t.keyPromptFlipFlop()
}

func (t *titleScene) update() bool {
	// Update bot snakes
	param.TeleportEnabled = titleSceneAlive
	for iSnake := 0; iSnake < numBotSnakes; iSnake++ {
		t.snakes[iSnake].Update(param.MouthAnimStartDistance)
	}

	// Update player snake
	param.TeleportEnabled = true
	t.playerSnake.Update(param.MouthAnimStartDistance)

	if titleSceneAlive {
		t.handleKeyPress()

	} else {
		// Update transition process to the next scene
		t.titleRectAlpha -= titleRectDissapearRate
		t.titleRectDrawOpts.Uniforms["Alpha"] = t.titleRectAlpha
		if t.titleRectAlpha <= 0.0 {
			t.shaderTitle.Dispose()
			return true
		}
	}

	return false
}

func (t *titleScene) handleKeyPress() {
	t.pressedKeys = inpututil.AppendPressedKeys(t.pressedKeys[:0])
	if len(t.pressedKeys) > 0 && titleSceneAlive {
		// Start transition process
		titleSceneAlive = false
		t.titleRectDrawOpts.Uniforms["ShowKeyPrompt"] = float32(0.0)

		// Increase speeds of snakes other than the player's snake
		for iSnake := 0; iSnake < numBotSnakes; iSnake++ {
			t.snakes[iSnake].Speed *= dumbSnakeRunMultip
		}
	}
}

func (t *titleScene) draw(screen *ebiten.Image) {
	screen.Fill(param.ColorBackground)

	// Draw bot snakes
	for iSnake := 0; iSnake < numBotSnakes; iSnake++ {
		t.snakes[iSnake].Draw(screen)
	}

	// Draw player snake
	t.playerSnake.Draw(screen)

	drawFPS(screen)

	// Draw Title Rect
	vertices, indices := t.titleRectComp.Triangles()
	screen.DrawTrianglesShader(vertices, indices, t.shaderTitle, &t.titleRectDrawOpts)
}

// Goroutine
// keyPromptFlipFlop switches title rectangle image within a specific time
func (t *titleScene) keyPromptFlipFlop() {
	const showTimeHalfSecs = keyPromptShowTimeSec * 2
	const hideTimeHalfSecs = keyPromptHideTimeSec * 2
	showPrompt := true

	halfSecondTicker := time.NewTicker(time.Millisecond * 500)
	for titleSceneAlive {
		if showPrompt {
			t.titleRectDrawOpts.Uniforms["ShowKeyPrompt"] = float32(1.0)
			for ihalfSecs := 0; ihalfSecs < showTimeHalfSecs; ihalfSecs++ {
				<-halfSecondTicker.C
			}
		} else {
			t.titleRectDrawOpts.Uniforms["ShowKeyPrompt"] = float32(0.0)
			for ihalfSecs := 0; ihalfSecs < hideTimeHalfSecs; ihalfSecs++ {
				<-halfSecondTicker.C
			}
		}
		showPrompt = !showPrompt
	}
	halfSecondTicker.Stop()
}

// Goroutine
// control turns given snake at a given time
func control(snake *s.Snake) {
	const turnTimeMinMs = turnTimeMinSec * 1000
	const turnTimeDiffMs = turnTimeMaxSec*1000 - turnTimeMinMs

	var dirNew s.DirectionT
	for titleSceneAlive {
		// Determine the new direction.
		dirCurrent := snake.LastDirection()

		if randResult := rand.Float32(); dirCurrent.IsVertical() {
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

		// Sleep a random amount of time between turnTimeMin and turnTimeMax.
		sleepTime := time.Duration(turnTimeMinMs + rand.Float32()*turnTimeDiffMs)
		time.Sleep(time.Millisecond * sleepTime)
	}
}
