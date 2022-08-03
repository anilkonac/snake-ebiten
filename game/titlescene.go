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
	titleRectDissapearRate float64 = (80.0 / 255.0) * param.DeltaTime
	textTitle                      = "Ssnake"
	textPressToPlay                = "Press any key to start"
	textTitleShiftY                = -50
	textKeyPromptShiftY            = +100
	keyPromptShowTimeSec           = 1.0
	keyPromptHideTimeSec           = 0.5
)

var (
	titleSceneAlive = true
	colorTitleRect  = []float64{float64(param.ColorSnake2.R) / 255.0, float64(param.ColorSnake2.G) / 255.0, float64(param.ColorSnake2.B) / 255.0, titleRectInitialAlpha}

	snakeColors = [...]*color.RGBA{
		&param.ColorSnake1,
		&param.ColorSnake2,
		&param.ColorFood,
		&param.ColorDebug,
	}
	rectShader ebiten.Shader
)

type titleScene struct {
	titleRectAlpha      float64
	turnTimers          [numSnakes]float32
	sceneTime           float32
	snakes              []*s.Snake
	pressedKeys         []ebiten.Key
	titleImage          *ebiten.Image
	titleImageKeyPrompt *ebiten.Image
	titleRectDrawOpts   ebiten.DrawImageOptions
}

func newTitleScene(playerSnake *s.Snake) *titleScene {

	// Create scene
	scene := &titleScene{
		titleRectAlpha: titleRectInitialAlpha,
		snakes:         make([]*s.Snake, numSnakes),
		pressedKeys:    make([]ebiten.Key, 0, 10),
	}
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
		scene.turnTimers[iSnake] = turnTimeMinSec + rand.Float32()*turnTimeDiff

	}

	// Store playersnake pointer as an element of the snakes array
	scene.snakes[numSnakes-1] = playerSnake

	// Set rand turn time
	scene.turnTimers[numSnakes-1] = turnTimeMinSec + rand.Float32()*turnTimeDiff

	return scene
}

func (t *titleScene) prepareTitleRects() {
	boundTextTitleSize := boundTextTitle.Size()
	boundTextKeyPromptSize := boundTextKeyPrompt.Size()

	rectShader = *shader.New(shader.PathTitleRect)

	// Prepare title rect image
	t.titleImage = ebiten.NewImage(titleRectWidth, titleRectHeight)
	textImage := ebiten.NewImage(titleRectWidth, titleRectHeight)

	// Draw Title text to the image
	text.Draw(textImage, textTitle, fontFaceTitle,
		(titleRectWidth-boundTextTitleSize.X)/2.0-boundTextTitle.Min.X,
		(titleRectHeight-boundTextTitleSize.Y)/2.0-boundTextTitle.Min.Y+textTitleShiftY,
		color.White)

	textImageWPrompt := ebiten.NewImageFromImage(textImage)

	t.titleImage.DrawRectShader(titleRectWidth, titleRectHeight, &rectShader, &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]interface{}{
			"RadiusTex": []float32{float32(titleRectCornerRadiusX / titleRectWidth), float32(titleRectCornerRadiusY / titleRectHeight)},
		},
		Images: [4]*ebiten.Image{
			textImage,
		},
	})

	// Prepare key prompt text image
	t.titleImageKeyPrompt = ebiten.NewImage(titleRectWidth, titleRectHeight)

	// Draw key prompt text to the image
	text.Draw(textImageWPrompt, textPressToPlay, param.FontFaceScore,
		(titleRectWidth-boundTextKeyPromptSize.X)/2.0-boundTextKeyPrompt.Min.X,
		(titleRectHeight-boundTextKeyPromptSize.Y)/2.0-boundTextKeyPrompt.Min.Y+textKeyPromptShiftY, color.White)

	t.titleImageKeyPrompt.DrawRectShader(titleRectWidth, titleRectHeight, &rectShader, &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]interface{}{
			"RadiusTex": []float32{float32(titleRectCornerRadiusX / titleRectWidth), float32(titleRectCornerRadiusY / titleRectHeight)},
		},
		Images: [4]*ebiten.Image{
			textImageWPrompt,
		},
	})

	// Initialize title rect draw options
	t.titleRectDrawOpts.GeoM.Translate((param.ScreenWidth-titleRectWidth)/2.0, (param.ScreenHeight-titleRectHeight)/2.0)
	t.titleRectDrawOpts.ColorM.Scale(colorTitleRect[0], colorTitleRect[1], colorTitleRect[2], colorTitleRect[3])

	// go t.keyPromptFlipFlop()
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
			colorTitleRect[3] = t.titleRectAlpha
			// TODO: Find more efficient way
			t.titleRectDrawOpts.ColorM.Reset()
			t.titleRectDrawOpts.ColorM.Scale(colorTitleRect[0], colorTitleRect[1], colorTitleRect[2], colorTitleRect[3])
			if t.titleRectAlpha <= 0.0 {
				// t.shaderTitle.Dispose()
				rectShader.Dispose()
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
		t.turnTimers[iSnake] += turnTimeMinSec + rand.Float32()*turnTimeDiff
	}

	snake.Update(param.MouthAnimStartDistance) // Make sure the snake's mouth is not open
}

func (t *titleScene) handleKeyPress() {
	t.pressedKeys = inpututil.AppendPressedKeys(t.pressedKeys[:0])
	if len(t.pressedKeys) > 0 && titleSceneAlive {
		// Start transition process
		titleSceneAlive = false
		// t.titleRectDrawOpts.Uniforms["ShowKeyPrompt"] = float32(0.0)

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
		snake.Draw(screen)
	}

	drawFPS(screen)

	// Draw Title Rect
	screen.DrawImage(t.titleImage, &t.titleRectDrawOpts)
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

// // Goroutine
// func (t *titleScene) keyPromptFlipFlop() {
// 	showTimeHalfSecs := int(keyPromptShowTime * 2)
// 	hideTimeHalfSecs := int(keyPromptHideTime * 2)
// 	showPrompt := true

// 	halfSecondTicker := time.NewTicker(time.Millisecond * 500)
// 	for titleSceneAlive {
// 		if showPrompt {
// 			t.titleRectDrawOpts.Uniforms["ShowKeyPrompt"] = float32(1.0)
// 			for ihalfSecs := 0; ihalfSecs < showTimeHalfSecs; ihalfSecs++ {
// 				<-halfSecondTicker.C
// 			}
// 		} else {
// 			t.titleRectDrawOpts.Uniforms["ShowKeyPrompt"] = float32(0.0)
// 			for ihalfSecs := 0; ihalfSecs < hideTimeHalfSecs; ihalfSecs++ {
// 				<-halfSecondTicker.C
// 			}
// 		}
// 		showPrompt = !showPrompt
// 	}
// 	halfSecondTicker.Stop()
// }
