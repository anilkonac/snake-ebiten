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
	"fmt"
	"image/color"

	"github.com/anilkonac/snake-ebiten/game/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Game constants
const (
	ScreenWidth  = 960
	ScreenHeight = 720
	deltaTime    = 1.0 / 60.0
	restartTime  = 1.5 // seconds
)

// Snake parameters
const (
	snakeHeadCenterX  = ScreenWidth / 2.0
	snakeHeadCenterY  = ScreenHeight / 2.0
	snakeSpeedInitial = 300
	snakeSpeedFinal   = 275
	snakeLength       = 240
	snakeWidth        = 30
)

const halfSnakeWidth = snakeWidth / 2.0

const (
	dpi             = 84
	fontSizeScore   = 24
	fontSizeDebug   = 16
	scoreTextShiftX = 10
	scoreTextShiftY = 8
	fpsTextShiftX   = 3
	fpsTextShiftY   = 2
)

// Colors to be used in the drawing.
// Palette: https://coolors.co/palette/003049-d62828-f77f00-fcbf49-eae2b7
var (
	colorBackground = color.RGBA{0, 48, 73, 255}     // ~ Prussian Blue
	colorSnake1     = color.RGBA{252, 191, 73, 255}  // ~ Maximum Yellow Red
	colorSnake2     = color.RGBA{247, 127, 0, 255}   // ~ Orange
	colorFood       = color.RGBA{214, 40, 40, 255}   // ~ Maximum Red
	colorDebug      = color.RGBA{234, 226, 183, 255} // ~ Lemon Meringue
	colorScore      = color.RGBA{247, 127, 0, 255}   // ~ Orange
)

var (
	printFPS   = true
	debugUnits = false // Draw consecutive units with different colors
	fontScore  font.Face
	fontDebug  font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.Rounded)
	if err != nil {
		panic(err)
	}

	fontScore, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeScore,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	tt, err = opentype.Parse(fonts.Debug)
	if err != nil {
		panic(err)
	}

	fontDebug, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeDebug,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
}

// Game implements ebiten.Game interface.
type Game struct {
	snake             *snake
	food              *food
	gameOver          bool
	paused            bool
	timeAfterGameOver float32
}

func NewGame() *Game {
	return &Game{
		snake: newSnake(snakeHeadCenterX, snakeHeadCenterY, directionRight, snakeLength),
		food:  newFoodRandLoc(),
	}
}

func (g *Game) restart() {
	*g = Game{
		snake: newSnakeRandDir(snakeHeadCenterX, snakeHeadCenterY, snakeLength),
		food:  newFoodRandLoc(),
	}
}

// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	g.handleSettingsInputs()

	if g.paused {
		return nil
	}

	repeatMusic()

	if g.gameOver {
		g.timeAfterGameOver += deltaTime
		if g.timeAfterGameOver >= restartTime {
			g.restart()
		}
		return nil
	}

	g.handleInput()
	g.snake.update()
	g.snake.checkIntersection(&g.gameOver)
	g.checkFood()

	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)

	// Draw food
	draw(screen, g.food)

	// Draw the snake
	for unit := g.snake.unitHead; unit != nil; unit = unit.next {
		draw(screen, unit)
	}

	// Draw score text
	g.drawScore(screen)

	if debugUnits {
		// Mark cursor
		x, y := ebiten.CursorPosition()
		markPoint(screen, float64(x), float64(y), 5, colorSnake2)

		// Print mouse coordinates
		msg := fmt.Sprintf("%d %d", x, y)
		rect := text.BoundString(fontDebug, msg)
		text.Draw(screen, msg, fontDebug, 0, -rect.Min.Y+ScreenHeight-rect.Size().Y, colorDebug)
	}

	g.printDebugMsgs(screen)
}

func (g *Game) drawScore(screen *ebiten.Image) {
	msg := fmt.Sprintf("Score: %d", g.snake.foodEaten)
	bound := text.BoundString(fontScore, msg)
	text.Draw(screen, msg, fontScore, scoreTextShiftX, -bound.Min.Y+scoreTextShiftY, colorScore)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) handleInput() {
	pressedLeft := inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	pressedRight := inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	pressedUp := inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
	pressedDown := inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)

	if !pressedLeft && !pressedRight && !pressedUp && !pressedDown {
		return
	}

	// Determine the new direction.
	dirCurrent := g.snake.lastDirection()
	dirNew := dirCurrent
	if dirCurrent.isVertical() {
		if pressedLeft {
			dirNew = directionLeft
		} else if pressedRight {
			dirNew = directionRight
		}
	} else {
		if pressedUp {
			dirNew = directionUp
		} else if pressedDown {
			dirNew = directionDown
		}
	}

	if dirNew == dirCurrent {
		return
	}

	// Create a new turn and take it
	newTurn := newTurn(dirCurrent, dirNew)
	g.snake.turnTo(newTurn, false)

}

func (g *Game) handleSettingsInputs() {
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		debugUnits = !debugUnits
		var numUnit uint8
		for unit := g.snake.unitHead; unit != nil; unit = unit.next {
			color := &colorSnake1
			if debugUnits && (numUnit%2 == 1) {
				color = &colorSnake2
			}
			unit.color = color
			numUnit++
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.paused = !g.paused
		if g.paused {
			playerMusic.Pause()
		} else {
			playerMusic.Play()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		printFPS = !printFPS
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		playMusic = !playMusic
		if playMusic {
			playerMusic.Play()
		} else {
			playerMusic.Pause()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		playSounds = !playSounds
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeyN) {
	// 	g.snake.grow()
	// 	g.snake.grow()
	// }
}

func (g *Game) checkFood() {
	if !g.food.isActive {
		// If food has spawned on the snake, respawn it elsewhere.
		for unit := g.snake.unitHead; unit != nil; unit = unit.next {
			if collides(unit, g.food, toleranceDefault) {
				g.food = newFoodRandLoc()
				return
			}
		}
		// Food has spawned in an open position, activate it.
		g.food.isActive = true
		return
	}

	if collides(g.snake.unitHead, g.food, toleranceFood) {
		g.snake.grow()
		g.food = newFoodRandLoc()
		playSoundEating()
		return
	}
}

func (g *Game) printDebugMsgs(screen *ebiten.Image) {
	if printFPS {
		msg := fmt.Sprintf("TPS: %.1f\tFPS: %.1f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
		bound := text.BoundString(fontDebug, msg)
		text.Draw(screen, msg, fontDebug, ScreenWidth-bound.Size().X-fpsTextShiftX, -bound.Min.Y+fpsTextShiftY, colorDebug)
	}
	// var totalLength float64
	// for unit := g.snake.unitHead; unit != nil; unit = unit.next {
	// 	totalLength += unit.length
	// }
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Food Eaten: %d   Snake length: %.2f   Speed: %.3f", g.snake.foodEaten, totalLength,  g.snake.speed))
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Turn Queue Length: %d Cap: %d", len(g.snake.turnQueue), cap(g.snake.turnQueue)), 0, 15)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Distance after turn: %.2f", g.snake.distAfterTurn), 0, 30)
}
