package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	snakeHeadCenterX = screenWidth / 2.0
	snakeHeadCenterY = screenHeight / 2.0
	snakeSpeed       = 200
	snakeDirection   = directionRight
	snakeLength      = 3
	unitLength       = 25
)

const deltaTime = 1.0 / 60.0

var (
	tps float64
	// mouseX int
	// mouseY int
)

// game implements ebiten.game interface.
type game struct {
	snake *snake
}

func newGame() *game {
	game := new(game)
	game.snake = newSnake(snakeHeadCenterX, snakeHeadCenterY, snakeDirection, snakeSpeed, snakeLength)

	return game
}

// Update is called every tick (1/60 [s] by default).
func (g *game) Update() error {
	tps = ebiten.CurrentTPS()
	// mouseX, mouseY = ebiten.CursorPosition()

	g.snake.update()

	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *game) Draw(screen *ebiten.Image) {
	g.printDebugMsgs(screen)

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
		curUnit.draw(screen, unitColor)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *game) printDebugMsgs(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.1f", tps))
	head := &g.snake.units[0]
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Head X: %.2f Y: %.2f", head.centerX, head.centerY), 0, 15)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse X: %d Y: %d", mouseX, mouseY), 0, 30)
}
