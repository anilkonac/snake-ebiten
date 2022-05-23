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
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// Game scene constants
const (
	restartTime      = 1.5 // seconds
	halfScreenWidth  = ScreenWidth / 2.0
	halfScreenHeight = ScreenHeight / 2.0
)

var (
	printFPS   = true
	debugUnits = false // Draw consecutive units with different colors
)

// gameScene implements ebiten.gameScene interface.
type gameScene struct {
	snake             *snake
	food              *food
	gameOver          bool
	paused            bool
	timeAfterGameOver float32
	scoreAnimList     []*scoreAnim
}

func newGameScene() *gameScene {
	teleportActive = true

	return &gameScene{
		snake: newSnake(vec64{snakeHeadCenterX, snakeHeadCenterY}, snakeLength, snakeSpeedInitial, directionRight, &colorSnake1),
		food:  newFoodRandLoc(),
	}
}

func (g *gameScene) restart() {
	*g = gameScene{
		snake: newSnakeRandDir(vec64{snakeHeadCenterX, snakeHeadCenterY}, snakeLength, snakeSpeedInitial, &colorSnake1),
		food:  newFoodRandLoc(),
	}
}

func (g *gameScene) update() bool {
	g.handleSettingsInputs()

	if g.paused {
		return false
	}

	if g.gameOver {
		g.timeAfterGameOver += deltaTime
		if g.timeAfterGameOver >= restartTime {
			g.restart()
		}
		return false
	}

	g.handleInput()

	distToFood := g.calcFoodDist()
	g.snake.update(distToFood)
	g.snake.checkIntersection(&g.gameOver)
	g.updateScoreAnims()
	g.checkFood(distToFood)

	return false
}

func (g *gameScene) calcFoodDist() float32 {
	if !g.food.isActive {
		return eatingAnimStartDistance
	}

	headLoc := g.snake.unitHead.headCenter
	foodLoc := g.food.center.to64()

	// In screen distance
	minDist := distance(headLoc, foodLoc)

	if headLoc.x < halfScreenWidth { // Left mirror distance
		mirroredFood := vec64{foodLoc.x - ScreenWidth, foodLoc.y}
		minDist = math.Min(minDist, distance(headLoc, mirroredFood))
	} else if headLoc.x >= halfScreenWidth { // Right mirror distance
		mirroredFood := vec64{foodLoc.x + ScreenWidth, foodLoc.y}
		minDist = math.Min(minDist, distance(headLoc, mirroredFood))
	}

	if headLoc.y < halfScreenHeight { // Upper mirror distance
		mirroredFood := vec64{foodLoc.x, foodLoc.y - ScreenHeight}
		minDist = math.Min(minDist, distance(headLoc, mirroredFood))
	} else if headLoc.y >= halfScreenHeight { // Bottom mirror distance
		mirroredFood := vec64{foodLoc.x, foodLoc.y + ScreenHeight}
		minDist = math.Min(minDist, distance(headLoc, mirroredFood))
	}

	return float32(minDist)

}

func (g *gameScene) updateScoreAnims() {
	for index, scoreAnim := range g.scoreAnimList {
		if scoreAnim.update() {
			g.scoreAnimList = append(g.scoreAnimList[:index], g.scoreAnimList[index+1:]...) // Delete score anim
			break
		}
	}
}

func (g *gameScene) handleInput() {
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

func (g *gameScene) handleSettingsInputs() {
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
		if g.paused && (musicState == musicOn) {
			playerMusic.Pause()
			musicState = musicPaused
		} else if musicState == musicPaused {
			playerMusic.Play()
			musicState = musicOn
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		if musicState == musicOn {
			musicState = musicMuted
			playerMusic.Pause()
		} else if (musicState == musicMuted) && !g.paused {
			musicState = musicOn
			playerMusic.Play()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		printFPS = !printFPS
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		playSounds = !playSounds
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeyN) {
	// 	g.snake.grow()
	// 	g.snake.grow()
	// }
}

func (g *gameScene) checkFood(distToFood float32) {
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

	if distToFood <= radiusEating {
		g.snake.grow()
		g.triggerScoreAnim()
		g.food = newFoodRandLoc()
		playSoundEating()
	}
}

func (g *gameScene) triggerScoreAnim() {
	corrCenter := g.snake.unitHead.headCenter

	// Correct the x and y position so the base score animation position will be the tip of the head,
	// not the head center.
	switch g.snake.unitHead.direction {
	case directionUp:
		corrCenter.y -= radiusSnake
	case directionDown:
		corrCenter.y += radiusSnake
	case directionRight:
		corrCenter.x += radiusSnake
	case directionLeft:
		corrCenter.x -= radiusSnake
	}

	g.scoreAnimList = append(g.scoreAnimList, newScoreAnim(corrCenter.to32()))
}

func (g *gameScene) draw(screen *ebiten.Image) {
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

	drawFPS(screen)

	if debugUnits {
		// Mark cursor
		x, y := ebiten.CursorPosition()
		markPoint(screen, vecI{x, y}.to64(), 5, colorSnake2)

		// Print mouse coordinates
		msg := fmt.Sprintf("%d %d", x, y)
		rect := text.BoundString(fontFaceDebug, msg)
		text.Draw(screen, msg, fontFaceDebug, 0, -rect.Min.Y+ScreenHeight-rect.Size().Y, colorDebug)
	}

	g.printDebugMsgs(screen)
}

func (g *gameScene) drawScore(screen *ebiten.Image) {
	msg := fmt.Sprintf("Score: %05d", int(g.snake.foodEaten)*foodScore)
	text.Draw(screen, msg, fontFaceScore, scoreTextShiftX, -boundTextScore.Min.Y+scoreTextShiftY, colorScore)
}

func (g *gameScene) printDebugMsgs(screen *ebiten.Image) {
	// var totalLength float64
	// for unit := g.snake.unitHead; unit != nil; unit = unit.next {
	// 	totalLength += unit.length
	// }
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Food Eaten: %d   Snake length: %.2f   Speed: %.3f", g.snake.foodEaten, totalLength,  g.snake.speed))
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Turn Queue Length: %d Cap: %d", len(g.snake.turnQueue), cap(g.snake.turnQueue)), 0, 15)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Distance after turn: %.2f", g.snake.distAfterTurn), 0, 30)
}
