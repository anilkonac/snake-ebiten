package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Snake parameters
const (
	snakeHeadCenterX = screenWidth / 2.0
	snakeHeadCenterY = screenHeight / 2.0
	snakeSpeed       = 200
	snakeDirection   = directionRight
	snakeLength      = 4
	unitLength       = 25
)

// Game constants
const (
	deltaTime      = 1.0 / 60.0
	halfUnitLength = unitLength / 2.0
)

// Colors to be used for drawing
var (
	colorBackground = color.RGBA{7, 59, 76, 255}     // Midnight Green Eagle Green
	colorSnake      = color.RGBA{255, 209, 102, 255} // Orange Yellow Crayola
)

// Debug variables
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

	g.handleInput()
	g.snake.update()

	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)
	g.snake.draw(screen)
	g.printDebugMsgs(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *game) printDebugMsgs(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.1f", tps))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Head X: %.2f Y: %.2f", g.snake.headCenterX, g.snake.headCenterY), 0, 15)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse X: %d Y: %d", mouseX, mouseY), 0, 30)
}

func (g *game) handleInput() {
	pressedLeft := inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	pressedRight := inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	pressedUp := inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
	pressedDown := inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)

	snakeDir := g.snake.direction
	if snakeDir == directionUp || snakeDir == directionDown {
		if pressedLeft {
			g.snake.rotateTo(directionLeft)
		} else if pressedRight {
			g.snake.rotateTo(directionRight)
		}
	} else if snakeDir == directionLeft || snakeDir == directionRight {
		if pressedUp {
			g.snake.rotateTo(directionUp)
		} else if pressedDown {
			g.snake.rotateTo(directionDown)
		}
	}
}
