package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	snakeHeadPosX  = screenWidth / 2.0
	snakeHeadPosY  = screenHeight / 2.0
	snakeSpeed     = 200
	snakeLength    = 5
	snakeDirection = directionRight
	deltaTime      = 1.0 / 60.0
)

// game implements ebiten.game interface.
type game struct {
	snake *snake
	tps   float64
	// mouseX int
	// mouseY int
}

func newGame() *game {
	game := new(game)
	game.snake = newSnake(snakeHeadPosX, snakeHeadPosY, snakeDirection, snakeSpeed, snakeLength)

	return game
}

// Update is called every tick (1/60 [s] by default).
func (g *game) Update() error {
	g.tps = ebiten.CurrentTPS()
	// g.mouseX, g.mouseY = ebiten.CursorPosition()

	// Update snake position
	headDirection := g.snake.headDirection()
	frameTravelDistance := snakeSpeed * deltaTime
	switch headDirection {
	case directionUp:
		g.snake.headPosY -= frameTravelDistance
	case directionDown:
		g.snake.headPosY += frameTravelDistance
	case directionRight:
		g.snake.headPosX += frameTravelDistance
	case directionLeft:
		g.snake.headPosX -= frameTravelDistance
	}

	// TODO: Update relative position and direction of units

	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *game) Draw(screen *ebiten.Image) {
	// Print TPS
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.1f", g.tps))
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse X: %d Y: %d", g.mouseX, g.mouseY), 0, 20)

	// Draw snake
	for indexUnit := 0; indexUnit < len(g.snake.units); indexUnit++ {
		// Get unit
		curUnit := &g.snake.units[indexUnit]

		// Compute position of current unit
		posX := g.snake.headPosX + float64(curUnit.relX)
		posY := g.snake.headPosY + float64(curUnit.relY)

		// Define color of the unit
		var unitColor color.Color
		if indexUnit%2 == 0 {
			unitColor = color.White
		} else {
			unitColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
		}

		// Draw unit
		ebitenutil.DrawRect(screen, posX, posY, unitLength, unitLength, unitColor)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
