package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

// Snake parameters
const (
	snakeHeadCenterX      = ScreenWidth / 2.0
	snakeHeadCenterY      = ScreenHeight / 2.0
	snakeSpeed            = 200
	snakeLength           = 200
	snakeWidth            = 25
	debugRects            = false // Draw consecutive units with different colors and draw position info of rects.
	lengthIncreasePercent = 18
)

// Game constants
const (
	deltaTime      = 1.0 / 60.0
	halfSnakeWidth = snakeWidth / 2.0
	restartTime    = 1.5 // seconds
)

// Colors to be used in the drawing.
var (
	colorBackground = color.RGBA{7, 59, 76, 255}     // Midnight Green Eagle Green
	colorSnake1     = color.RGBA{255, 209, 102, 255} // Orange Yellow Crayola
	colorSnake2     = color.RGBA{6, 214, 160, 255}   // Caribbean Green
	colorFood       = color.RGBA{239, 71, 111, 255}  // Paradise Pink
)

// Game implements ebiten.Game interface.
type Game struct {
	snake             *snake
	food              *food
	gameOver          bool
	paused            bool
	timeAfterGameOver float32
}

func NewGame() *Game {
	return &Game{
		snake: newSnake(snakeHeadCenterX, snakeHeadCenterY, directionRight, snakeSpeed, snakeLength),
		food:  newFoodRandLoc(),
	}
}

func (g *Game) restart() {
	*g = Game{
		snake: newSnakeRandDir(snakeHeadCenterX, snakeHeadCenterY, snakeSpeed, snakeLength),
		food:  newFoodRandLoc(),
	}
}

// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.paused = !g.paused
	}

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
	g.snake.update()
	g.gameOver = g.snake.checkIntersection()
	g.checkFood()

	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)

	g.snake.draw(screen)
	draw(g.food, screen)

	g.printDebugMsgs(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
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
	if isVertical(dirCurrent) {
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

func (g *Game) checkFood() {
	if !g.food.active {
		// If food has spawned on the snake, respawn it elsewhere.
		for curUnit := g.snake.unitHead; curUnit != nil; curUnit = curUnit.next {
			if collides(curUnit, g.food) {
				g.food = newFoodRandLoc()
				return
			}
		}
		// Food has spawned in an open position, activate it.
		g.food.active = true
		return
	}

	// Check if the snake has eaten the food.
	// for _, rectHead := range g.snake.unitHead.rects {
	// 	for _, rectFood := range g.food.rects {
	// 		if !intersects(rectHead, rectFood) {
	// 			continue
	// 		}

	// 		g.snake.grow()
	// 		g.food = newFoodRandLoc()
	// 		return
	// 	}
	// }
	if collides(g.snake.unitHead, g.food) {
		g.snake.grow()
		g.food = newFoodRandLoc()
		return
	}
}

func (g *Game) printDebugMsgs(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.1f  FPS: %.1f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
	// headUnit := g.snake.unitHead
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Head X: %.2f Y: %.2f", headUnit.headCenterX, headUnit.headCenterY), 0, 15)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse X: %d Y: %d", mouseX, mouseY), 0, 30)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Turn Queue Length: %d Cap: %d", len(g.snake.turnQueue), cap(g.snake.turnQueue)), 0, 15)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Distance after turn: %.2f", g.snake.distAfterTurn), 0, 30)
}
