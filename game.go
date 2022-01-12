package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	snakeHeadCenterX = screenWidth - 15
	snakeHeadCenterY = screenHeight / 2.0
	snakeSpeed       = 1
	snakeLength      = 5
	deltaTime        = 1.0 / 60.0
)

// Game implements ebiten.Game interface.
type Game struct {
	tps   float64
	snake *snake
}

func newGame() *Game {
	game := new(Game)
	game.snake = newSnake(snakeHeadCenterX, snakeHeadCenterY, directionRight, snakeSpeed, snakeLength)

	return game
}

// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	g.tps = ebiten.CurrentTPS()

	g.snake.update()

	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Print TPS
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.1f", g.tps))
	head := &g.snake.units[0]
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Head X: %.2f Y: %.2f", head.posX, head.posY), 0, 15)

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
		ebitenutil.DrawRect(screen, curUnit.posX-centerOffset, curUnit.posY-centerOffset, unitLength, unitLength, unitColor)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
