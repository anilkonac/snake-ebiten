package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	snakeHeadPosX = 200
	snakeHeadPosY = 200
	snakeSpeed    = 200
	snakeLength   = 5
	deltaTime     = 1.0 / 60.0
)

// Game implements ebiten.Game interface.
type Game struct {
	tps   float64
	snake *Snake
}

func newGame() *Game {
	game := new(Game)
	game.snake = newSnake(snakeHeadPosX, snakeHeadPosY, directionRight, snakeSpeed, snakeLength)

	return game
}

// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	g.tps = ebiten.CurrentTPS()

	// Update snake position
	for indexUnit := 0; indexUnit < len(g.snake.units); indexUnit++ {
		curUnit := &g.snake.units[indexUnit]

		travelDistance := float64(g.snake.speed) * deltaTime
		switch curUnit.direction {
		case directionRight:
			curUnit.posX += travelDistance
		case directionLeft:
			curUnit.posX -= travelDistance
		case directionUp:
			curUnit.posY -= travelDistance
		case directionDown:
			curUnit.posY += travelDistance
		}
	}

	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Print TPS
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.1f", g.tps))

	// Draw snake
	for indexUnit := 0; indexUnit < len(g.snake.units); indexUnit++ {
		// Get unit
		curUnit := &g.snake.units[indexUnit]

		// Define color of the unit
		var unitColor color.Color
		if indexUnit%2 == 0 {
			unitColor = color.White
		} else {
			unitColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
		}

		// Draw unit
		ebitenutil.DrawRect(screen, curUnit.posX, curUnit.posY, unitLength, unitLength, unitColor)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
