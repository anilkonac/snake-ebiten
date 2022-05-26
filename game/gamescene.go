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

	"github.com/anilkonac/snake-ebiten/game/params"
	s "github.com/anilkonac/snake-ebiten/game/snake"
	t "github.com/anilkonac/snake-ebiten/game/tools"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// Game scene constants
const (
	restartTime = 1.5 // seconds

	snakeHeadCenterX = params.HalfScreenWidth
	snakeHeadCenterY = params.HalfScreenHeight
)

type gameScene struct {
	snake             *s.Snake
	food              *food
	gameOver          bool
	paused            bool
	timeAfterGameOver float32
	scoreAnimList     []*scoreAnim
}

func newGameScene() *gameScene {
	params.TeleportActive = true

	return &gameScene{
		snake: s.NewSnake(t.Vec64{X: snakeHeadCenterX, Y: snakeHeadCenterY}, params.SnakeLength, params.SnakeSpeedInitial, s.DirectionRight, &params.ColorSnake1),
		food:  newFoodRandLoc(),
	}
}

func (g *gameScene) restart() {
	*g = gameScene{
		snake: s.NewSnakeRandDir(t.Vec64{X: snakeHeadCenterX, Y: snakeHeadCenterY}, params.SnakeLength, params.SnakeSpeedInitial, &params.ColorSnake1),
		food:  newFoodRandLoc(),
	}
}

func (g *gameScene) update() bool {
	g.handleSettingsInputs()

	if g.paused {
		return false
	}

	if g.gameOver {
		g.timeAfterGameOver += params.DeltaTime
		if g.timeAfterGameOver >= restartTime {
			g.restart()
		}
		return false
	}

	g.handleInput()

	distToFood := g.calcFoodDist()
	g.snake.Update(distToFood)
	g.checkIntersection()
	g.updateScoreAnims()
	g.checkFood(distToFood)

	return false
}

func (g *gameScene) checkIntersection() {
	curUnit := g.snake.UnitHead.Next
	if curUnit == nil {
		return
	}

	var tolerance float32 = params.ToleranceDefault
	if len(curUnit.RectsCollision) > 1 { // If second unit is on an edge
		tolerance = params.ToleranceScreenEdge // To avoid false collisions on screen edges
	}

	for curUnit != nil {
		if collides(g.snake.UnitHead, curUnit, tolerance) {
			g.gameOver = true
			playSoundHit()
			return
		}
		curUnit = curUnit.Next
	}
}

func (g *gameScene) calcFoodDist() float32 {
	if !g.food.isActive {
		return params.EatingAnimStartDistance
	}

	headLoc := g.snake.UnitHead.HeadCenter
	foodLoc := g.food.center.To64()

	// In screen distance
	minDist := t.Distance(headLoc, foodLoc)

	if headLoc.X < params.HalfScreenWidth { // Left mirror distance
		mirroredFood := t.Vec64{X: foodLoc.X - params.ScreenWidth, Y: foodLoc.Y}
		minDist = math.Min(minDist, t.Distance(headLoc, mirroredFood))
	} else if headLoc.X >= params.HalfScreenWidth { // Right mirror distance
		mirroredFood := t.Vec64{X: foodLoc.X + params.ScreenWidth, Y: foodLoc.Y}
		minDist = math.Min(minDist, t.Distance(headLoc, mirroredFood))
	}

	if headLoc.Y < params.HalfScreenHeight { // Upper mirror distance
		mirroredFood := t.Vec64{X: foodLoc.X, Y: foodLoc.Y - params.ScreenHeight}
		minDist = math.Min(minDist, t.Distance(headLoc, mirroredFood))
	} else if headLoc.Y >= params.HalfScreenHeight { // Bottom mirror distance
		mirroredFood := t.Vec64{X: foodLoc.X, Y: foodLoc.Y + params.ScreenHeight}
		minDist = math.Min(minDist, t.Distance(headLoc, mirroredFood))
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
	dirCurrent := g.snake.LastDirection()
	dirNew := dirCurrent
	if dirCurrent.IsVertical() {
		if pressedLeft {
			dirNew = s.DirectionLeft
		} else if pressedRight {
			dirNew = s.DirectionRight
		}
	} else {
		if pressedUp {
			dirNew = s.DirectionUp
		} else if pressedDown {
			dirNew = s.DirectionDown
		}
	}

	if dirNew == dirCurrent {
		return
	}

	// Create a new turn and take it
	newTurn := s.NewTurn(dirCurrent, dirNew)
	g.snake.TurnTo(newTurn, false)

}

func (g *gameScene) handleSettingsInputs() {
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		params.DebugUnits = !params.DebugUnits
		var numUnit uint8
		for unit := g.snake.UnitHead; unit != nil; unit = unit.Next {
			color := &params.ColorSnake1
			if params.DebugUnits && (numUnit%2 == 1) {
				color = &params.ColorSnake2
			}
			unit.SetColor(color)
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
		params.PrintFPS = !params.PrintFPS
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		playSounds = !playSounds
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeyN) {
	// 	g.snake.Grow()
	// 	g.snake.Grow()
	// }
}

func (g *gameScene) checkFood(distToFood float32) {
	if !g.food.isActive {
		// If food has spawned on the snake, respawn it elsewhere.
		for unit := g.snake.UnitHead; unit != nil; unit = unit.Next {
			if collides(unit, g.food, params.ToleranceDefault) {
				g.food = newFoodRandLoc()
				return
			}
		}
		// Food has spawned in an open position, activate it.
		g.food.isActive = true
		return
	}

	if distToFood <= params.RadiusEating {
		g.snake.Grow()
		g.triggerScoreAnim()
		g.food = newFoodRandLoc()
		playSoundEating()
	}
}

func (g *gameScene) triggerScoreAnim() {
	corrCenter := g.snake.UnitHead.HeadCenter

	// Correct the x and y position so the base score animation position will be the tip of the head,
	// not the head center.
	switch g.snake.UnitHead.Direction {
	case s.DirectionUp:
		corrCenter.Y -= params.RadiusSnake
	case s.DirectionDown:
		corrCenter.Y += params.RadiusSnake
	case s.DirectionRight:
		corrCenter.X += params.RadiusSnake
	case s.DirectionLeft:
		corrCenter.X -= params.RadiusSnake
	}

	g.scoreAnimList = append(g.scoreAnimList, newScoreAnim(corrCenter.To32()))
}

func (g *gameScene) draw(screen *ebiten.Image) {
	screen.Fill(params.ColorBackground)

	// Draw food
	draw(screen, g.food)

	// Draw the snake
	for unit := g.snake.UnitHead; unit != nil; unit = unit.Next {
		draw(screen, unit)
	}

	// Draw score anim
	for _, scoreAnim := range g.scoreAnimList {
		draw(screen, scoreAnim)
	}

	// Draw score text
	g.drawScore(screen)

	drawFPS(screen)

	if params.DebugUnits {
		// Mark cursor
		x, y := ebiten.CursorPosition()
		t.MarkPoint(screen, t.VecI{X: x, Y: y}.To64(), 5, params.ColorSnake2)

		// Print mouse coordinates
		msg := fmt.Sprintf("%d %d", x, y)
		rect := text.BoundString(fontFaceDebug, msg)
		text.Draw(screen, msg, fontFaceDebug, 0, -rect.Min.Y+params.ScreenHeight-rect.Size().Y, params.ColorDebug)
	}

	g.printDebugMsgs(screen)
}

func (g *gameScene) drawScore(screen *ebiten.Image) {
	msg := fmt.Sprintf("Score: %05d", int(g.snake.FoodEaten)*params.FoodScore)
	text.Draw(screen, msg, fontFaceScore, scoreTextShiftX, -boundTextScore.Min.Y+scoreTextShiftY, params.ColorScore)
}

func (g *gameScene) printDebugMsgs(screen *ebiten.Image) {
	// var totalLength float64
	// for unit := g.snake.UnitHead; unit != nil; unit = unit.Next {
	// 	totalLength += unit.length
	// }
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Food Eaten: %d   Snake length: %.2f   Speed: %.3f", g.snake.foodEaten, totalLength,  g.snake.speed))
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Turn Queue Length: %d Cap: %d", len(g.snake.turnQueue), cap(g.snake.turnQueue)), 0, 15)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("t.Distance after turn: %.2f", g.snake.distAfterTurn), 0, 30)
}
