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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
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
	snakeHeadCenterX        = ScreenWidth / 2.0
	snakeHeadCenterY        = ScreenHeight / 2.0
	snakeSpeedInitial       = 300
	snakeSpeedFinal         = 275
	snakeLength             = 240
	snakeWidth              = 30
	eatingAnimStartDistance = 140
)

const halfSnakeWidth = snakeWidth / 2.0

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
)

// Game implements ebiten.Game interface.
type Game struct {
	snake             *snake
	food              *food
	gameOver          bool
	paused            bool
	timeAfterGameOver float32
	scoreAnimList     []*scoreAnim
}

func NewGame() *Game {
	return &Game{
		snake: newSnake(vec64{snakeHeadCenterX, snakeHeadCenterY}, directionRight, snakeLength),
		food:  newFoodRandLoc(),
	}
}

func (g *Game) restart() {
	*g = Game{
		snake: newSnakeRandDir(vec64{snakeHeadCenterX, snakeHeadCenterY}, snakeLength),
		food:  newFoodRandLoc(),
	}
}

// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	g.handleSettingsInputs()

	if g.paused {
		return nil
	}

	if g.gameOver {
		g.timeAfterGameOver += deltaTime
		if g.timeAfterGameOver >= restartTime {
			g.restart()
		}
		return nil
	}

	g.handleInput()

	// Calculate the distance betw head and food
	var distToFood float32 = eatingAnimStartDistance
	if g.food.isActive {
		distToFood = float32(distance(&g.snake.unitHead.headCenter, g.food.center.to64()))
	}

	g.snake.update(distToFood)
	g.snake.checkIntersection(&g.gameOver)
	g.updateScoreAnims()
	g.checkFood()

	return nil
}

func (g *Game) updateScoreAnims() {
	for index, scoreAnim := range g.scoreAnimList {
		if scoreAnim.update() {
			g.scoreAnimList = append(g.scoreAnimList[:index], g.scoreAnimList[index+1:]...) // Delete score anim
			break
		}
	}
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
			playMusic = false
		} else {
			playerMusic.Play()
			playMusic = true
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
		g.triggerScoreAnim()
		g.food = newFoodRandLoc()
		playSoundEating()
		return
	}
}

func (g *Game) triggerScoreAnim() {
	corrCenter := g.snake.unitHead.headCenter

	// Correct the x and y position so the base score animation position will be the tip of the head,
	// not the head center.
	switch g.snake.unitHead.direction {
	case directionUp:
		corrCenter.y -= halfSnakeWidth
	case directionDown:
		corrCenter.y += halfSnakeWidth
	case directionRight:
		corrCenter.x += halfSnakeWidth
	case directionLeft:
		corrCenter.x -= halfSnakeWidth
	}

	g.scoreAnimList = append(g.scoreAnimList, newScoreAnim(corrCenter.to32()))
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

	// Draw score anim
	for _, scoreAnim := range g.scoreAnimList {
		draw(screen, scoreAnim)
	}

	// Draw score text
	g.drawScore(screen)

	if debugUnits {
		// Mark cursor
		x, y := ebiten.CursorPosition()
		markPoint(screen, vecI{x, y}.to64(), 5, colorSnake2)

		// Print mouse coordinates
		msg := fmt.Sprintf("%d %d", x, y)
		rect := text.BoundString(fontDebug, msg)
		text.Draw(screen, msg, fontDebug, 0, -rect.Min.Y+ScreenHeight-rect.Size().Y, colorDebug)
	}

	g.printDebugMsgs(screen)
}

func (g *Game) drawScore(screen *ebiten.Image) {
	msg := fmt.Sprintf("Score: %05d", int(g.snake.foodEaten)*foodScore)
	text.Draw(screen, msg, fontScore, scoreTextShiftX, -boundScoreText.Min.Y+scoreTextShiftY, colorScore)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) printDebugMsgs(screen *ebiten.Image) {
	if printFPS {
		msg := fmt.Sprintf("TPS: %.1f\tFPS: %.1f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
		text.Draw(screen, msg, fontDebug, ScreenWidth-boundFPSText.Size().X-fpsTextShiftX, -boundFPSText.Min.Y+fpsTextShiftY, colorDebug)
	}
	// var totalLength float64
	// for unit := g.snake.unitHead; unit != nil; unit = unit.next {
	// 	totalLength += unit.length
	// }
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Food Eaten: %d   Snake length: %.2f   Speed: %.3f", g.snake.foodEaten, totalLength,  g.snake.speed))
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Turn Queue Length: %d Cap: %d", len(g.snake.turnQueue), cap(g.snake.turnQueue)), 0, 15)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Distance after turn: %.2f", g.snake.distAfterTurn), 0, 30)
}
