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

	"github.com/anilkonac/snake-ebiten/game/object"
	s "github.com/anilkonac/snake-ebiten/game/object/snake"
	"github.com/anilkonac/snake-ebiten/game/param"
	t "github.com/anilkonac/snake-ebiten/game/tool"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// Game scene constants
const (
	restartTime = 1.5 // seconds

	snakeHeadCenterX = param.HalfScreenWidth
	snakeHeadCenterY = param.HalfScreenHeight
)

type gameScene struct {
	snake             *s.Snake
	food              *object.Food
	gameOver          bool
	paused            bool
	timeAfterGameOver float32
	scoreAnimList     []*object.ScoreAnim
}

func newGameScene(snake *s.Snake) *gameScene {
	param.TeleportActive = true

	return &gameScene{
		snake: snake,
		food:  object.NewFoodRandLoc(),
	}
}

func (g *gameScene) restart() {
	*g = gameScene{
		snake: s.NewSnakeRandDir(t.Vec64{X: snakeHeadCenterX, Y: snakeHeadCenterY}, param.SnakeLength, param.SnakeSpeedInitial, &param.ColorSnake1),
		food:  object.NewFoodRandLoc(),
	}
}

func (g *gameScene) update() bool {
	g.handleSettingsInputs()

	if g.paused {
		return false
	}

	if g.gameOver {
		g.timeAfterGameOver += param.DeltaTime
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

	var tolerance float32 = param.ToleranceDefault
	if len(curUnit.RectsCollision) > 1 { // If second unit is on an edge
		tolerance = param.ToleranceScreenEdge // To avoid false collisions on screen edges
	}

	for curUnit != nil {
		if object.Collides(g.snake.UnitHead, curUnit, tolerance) {
			g.gameOver = true
			playSoundHit()
			return
		}
		curUnit = curUnit.Next
	}
}

func (g *gameScene) calcFoodDist() float32 {
	if !g.food.IsActive {
		return param.EatingAnimStartDistance
	}

	headLoc := g.snake.UnitHead.HeadCenter
	foodLoc := g.food.Center.To64()

	// In screen distance
	minDist := t.Distance(headLoc, foodLoc)

	if headLoc.X < param.HalfScreenWidth { // Left mirror distance
		mirroredFood := t.Vec64{X: foodLoc.X - param.ScreenWidth, Y: foodLoc.Y}
		minDist = math.Min(minDist, t.Distance(headLoc, mirroredFood))
	} else if headLoc.X >= param.HalfScreenWidth { // Right mirror distance
		mirroredFood := t.Vec64{X: foodLoc.X + param.ScreenWidth, Y: foodLoc.Y}
		minDist = math.Min(minDist, t.Distance(headLoc, mirroredFood))
	}

	if headLoc.Y < param.HalfScreenHeight { // Upper mirror distance
		mirroredFood := t.Vec64{X: foodLoc.X, Y: foodLoc.Y - param.ScreenHeight}
		minDist = math.Min(minDist, t.Distance(headLoc, mirroredFood))
	} else if headLoc.Y >= param.HalfScreenHeight { // Bottom mirror distance
		mirroredFood := t.Vec64{X: foodLoc.X, Y: foodLoc.Y + param.ScreenHeight}
		minDist = math.Min(minDist, t.Distance(headLoc, mirroredFood))
	}

	return float32(minDist)

}

func (g *gameScene) updateScoreAnims() {
	for index, scoreAnim := range g.scoreAnimList {
		if scoreAnim.Update() {
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
		param.DebugUnits = !param.DebugUnits
		var numUnit uint8
		for unit := g.snake.UnitHead; unit != nil; unit = unit.Next {
			color := &param.ColorSnake1
			if param.DebugUnits && (numUnit%2 == 1) {
				color = &param.ColorSnake2
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
		param.PrintFPS = !param.PrintFPS
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
	if !g.food.IsActive {
		// If food has spawned on the snake, respawn it elsewhere.
		for unit := g.snake.UnitHead; unit != nil; unit = unit.Next {
			if object.Collides(unit, g.food, param.ToleranceDefault) {
				g.food = object.NewFoodRandLoc()
				return
			}
		}
		// Food has spawned in an open position, activate it.
		g.food.IsActive = true
		return
	}

	// Check for collision with food
	if distToFood <= param.RadiusEating {
		g.snake.Grow()
		g.triggerScoreAnim()
		g.food = object.NewFoodRandLoc()
		playSoundEating()
	}
}

func (g *gameScene) triggerScoreAnim() {
	corrCenter := g.snake.UnitHead.HeadCenter

	// Correct the x and y position so the base score animation position will be the tip of the head,
	// not the head center.
	switch g.snake.UnitHead.Direction {
	case s.DirectionUp:
		corrCenter.Y -= param.RadiusSnake
	case s.DirectionDown:
		corrCenter.Y += param.RadiusSnake
	case s.DirectionRight:
		corrCenter.X += param.RadiusSnake
	case s.DirectionLeft:
		corrCenter.X -= param.RadiusSnake
	}

	g.scoreAnimList = append(g.scoreAnimList, object.NewScoreAnim(corrCenter.To32()))
}

func (g *gameScene) draw(screen *ebiten.Image) {
	screen.Fill(param.ColorBackground)

	// Draw food
	object.Draw(screen, g.food)

	// Draw the snake
	for unit := g.snake.UnitHead; unit != nil; unit = unit.Next {
		object.Draw(screen, unit)
	}

	// Draw score anim
	for _, scoreAnim := range g.scoreAnimList {
		object.Draw(screen, scoreAnim)
	}

	// Draw score text
	g.drawScore(screen)

	drawFPS(screen)

	if param.DebugUnits {
		// Mark cursor
		x, y := ebiten.CursorPosition()
		t.MarkPoint(screen, t.VecI{X: x, Y: y}.To64(), 5, param.ColorSnake2)

		// Print mouse coordinates
		msg := fmt.Sprintf("%d %d", x, y)
		rect := text.BoundString(fontFaceDebug, msg)
		text.Draw(screen, msg, fontFaceDebug, 0, -rect.Min.Y+param.ScreenHeight-rect.Size().Y, param.ColorDebug)
	}

	g.printDebugMsgs(screen)
}

func (g *gameScene) drawScore(screen *ebiten.Image) {
	msg := fmt.Sprintf("Score: %05d", int(g.snake.FoodEaten)*param.FoodScore)
	text.Draw(screen, msg, param.FontFaceScore, scoreTextShiftX, -boundTextScore.Min.Y+scoreTextShiftY, param.ColorScore)
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
