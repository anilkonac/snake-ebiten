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
	snakeLength      = 400
	snakeWidth       = 25
	debugUnits       = false // Draw consecutive units with different colors.
)

// Game constants
const (
	deltaTime      = 1.0 / 60.0
	halfSnakeWidth = snakeWidth / 2.0
	gameOverTime   = 1.5 // seconds
)

// Colors to be used for drawing
var (
	colorBackground = color.RGBA{7, 59, 76, 255}     // Midnight Green Eagle Green
	colorSnake1     = color.RGBA{255, 209, 102, 255} // Orange Yellow Crayola
	colorSnake2     = color.RGBA{239, 71, 111, 255}  // Paradise Pink
)

// Debug variables
var (
	tps float64
	// fps float64
)

// game implements ebiten.game interface.
type game struct {
	snake          *snake
	gameOver       bool
	paused         bool
	timeInGameOver float32
}

func newGame() *game {
	game := new(game)
	game.snake = newSnake(snakeHeadCenterX, snakeHeadCenterY, directionRight, snakeSpeed, snakeLength)

	return game
}

func (g *game) restart() {
	*g = game{
		snake: newSnakeRandDir(snakeHeadCenterX, snakeHeadCenterY, snakeSpeed, snakeLength),
	}
}

// Update is called every tick (1/60 [s] by default).
func (g *game) Update() error {
	tps = ebiten.CurrentTPS()
	// fps = ebiten.CurrentFPS()

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.paused = !g.paused
	}

	if g.paused {
		return nil
	}

	if g.gameOver {
		g.timeInGameOver += deltaTime
		if g.timeInGameOver >= gameOverTime {
			g.restart()
		}
		return nil
	}

	g.handleInput()
	g.snake.update()
	g.gameOver = g.snake.checkIntersection()

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

func (g *game) handleInput() {
	pressedLeft := inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	pressedRight := inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	pressedUp := inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
	pressedDown := inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)

	if !pressedLeft && !pressedRight && !pressedUp && !pressedDown {
		return
	}

	var dirCurrent directionT
	if queueLength := len(g.snake.turnQueue); queueLength > 0 {
		// Set current direction as the direction of the last turn to be taken.
		dirCurrent = g.snake.turnQueue[queueLength-1].directionTo
	} else {
		// Set as current head direction
		dirCurrent = g.snake.unitHead.direction
	}

	// Determine new direction
	dirNew := dirCurrent
	if dirCurrent == directionUp || dirCurrent == directionDown {
		if pressedLeft {
			dirNew = directionLeft
		} else if pressedRight {
			dirNew = directionRight
		}
	} else if dirCurrent == directionLeft || dirCurrent == directionRight {
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

func (g *game) printDebugMsgs(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.1f", tps))
	// headUnit := g.snake.unitHead
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Head X: %.2f Y: %.2f", headUnit.headCenterX, headUnit.headCenterY), 0, 15)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse X: %d Y: %d", mouseX, mouseY), 0, 30)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Turn Queue Length: %d Cap: %d", len(g.snake.turnQueue), cap(g.snake.turnQueue)), 0, 15)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Distance after turn: %.2f", g.snake.distAfterTurn), 0, 30)
}
