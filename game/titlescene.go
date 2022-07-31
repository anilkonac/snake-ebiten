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
	"github.com/anilkonac/snake-ebiten/game/object"
	s "github.com/anilkonac/snake-ebiten/game/object/snake"
	"github.com/anilkonac/snake-ebiten/game/param"
	"github.com/anilkonac/snake-ebiten/game/shader"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// Dumb snake parameters
const (
	numSnakes           = 30
	turnTimeMin         = 0 // sec
	turnTimeMax         = 2 // sec
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
	titleRectCornerRadiusX         = param.RadiusSnake
	titleRectCornerRadiusY         = titleRectCornerRadiusX / titleRectRatio
	titleRectInitialAlpha          = 230 / 255.0
	titleRectDissapearRate float32 = (80 / 255.0) * param.DeltaTime
	textTitle                      = "Ssnake"
	textPressToPlay                = "Press any key to start"
	textTitleShiftY                = -50
	textKeyPromptShiftY            = +100
	keyPromptShowTime              = 1.0 //sec
	keyPromptHideTime              = 0.5 // sec
)

var (
	titleSceneAlive = true
	colorTitleRect  = &param.ColorSnake2

	snakeColors = map[int]*color.RGBA{
		0: &param.ColorSnake1,
		1: &param.ColorSnake2,
		2: &param.ColorFood,
		3: &param.ColorDebug,
	}
)

type titleScene struct {
	titleRectComp     c.TeleCompTriang
	titleRectAlpha    float32
	turnTimers        [numSnakes]float32
	sceneTime         float32
	snakes            []*s.Snake
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
		titleRectAlpha: titleRectInitialAlpha,
		snakes:         make([]*s.Snake, numSnakes),
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
	for iSnake := 0; iSnake < numSnakes-1; iSnake++ {
		length := dumbSnakeLengthMin + rand.Intn(dumbSnakeLengthDiff)
		speed := dumbSnakeSpeedMin + rand.Float64()*dumbSnakeSpeedDiff
		snake := s.NewSnakeRandDirLoc(uint16(length), speed, snakeColors[rand.Intn(lenSnakeColors)])
		scene.snakes[iSnake] = snake

		// Set rand turn time
		scene.turnTimers[iSnake] = turnTimeMin + rand.Float32()*turnTimeDiff

	}

	// Store playersnake pointer as an element of the snakes array
	scene.snakes[numSnakes-1] = playerSnake

	// Set rand turn time
	scene.turnTimers[numSnakes-1] = turnTimeMin + rand.Float32()*turnTimeDiff

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
	if titleSceneAlive {
		for iSnake := range t.snakes {
			t.updateSnake(iSnake)
		}

		t.handleKeyPress()

	} else {
		// Update dumb snakes
		param.TeleportEnabled = false
		for iSnake := 0; iSnake < numSnakes-1; iSnake++ {
			t.updateSnake(iSnake)

		}

		// Update player snake
		param.TeleportEnabled = true
		t.updateSnake(numSnakes - 1)

		// Update transition process to the next scene
		if !titleSceneAlive {
			t.titleRectAlpha -= titleRectDissapearRate
			t.titleRectDrawOpts.Uniforms["Alpha"] = t.titleRectAlpha
			if t.titleRectAlpha <= 0.0 {
				t.shaderTitle.Dispose()
				return true
			}
		}
	}

	t.sceneTime += param.DeltaTime
	return false
}

func (t *titleScene) updateSnake(iSnake int) {
	snake := t.snakes[iSnake]
	if titleSceneAlive && t.sceneTime >= t.turnTimers[iSnake] {
		turnRandomly(snake)
		t.turnTimers[iSnake] += turnTimeMin + rand.Float32()*turnTimeDiff
	}

	snake.Update(param.MouthAnimStartDistance) // Make sure the snake's mouth is not open
}

func (t *titleScene) handleKeyPress() {
	t.pressedKeys = inpututil.AppendPressedKeys(t.pressedKeys[:0])
	if len(t.pressedKeys) > 0 && titleSceneAlive {
		// Start transition process
		titleSceneAlive = false
		t.titleRectDrawOpts.Uniforms["ShowKeyPrompt"] = float32(0.0)

		// Increase speeds of snakes other than the player's snake
		for iSnake := 0; iSnake < numSnakes-1; iSnake++ {
			t.snakes[iSnake].Speed *= dumbSnakeRunMultip
		}
	}
}

func (t *titleScene) draw(screen *ebiten.Image) {
	screen.Fill(param.ColorBackground)

	// Draw snakes
	for _, snake := range t.snakes {
		for unit := snake.UnitHead; unit != nil; unit = unit.Next {
			object.Draw(screen, unit)
		}

	}

	drawFPS(screen)

	// Draw Title Rect
	vertices, indices := t.titleRectComp.Triangles()
	screen.DrawTrianglesShader(vertices, indices, t.shaderTitle, &t.titleRectDrawOpts)
}

func turnRandomly(snake *s.Snake) {
	// Determine the new direction.
	var dirNew s.DirectionT
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
}

// Goroutine
func (t *titleScene) keyPromptFlipFlop() {
	showTimeHalfSecs := int(keyPromptShowTime * 2)
	hideTimeHalfSecs := int(keyPromptHideTime * 2)
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
