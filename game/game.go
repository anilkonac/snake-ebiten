package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 960
	ScreenHeight = 720
	deltaTime    = 1.0 / 60.0
)

var teleportActive = true

// Colors to be used in the drawing.
// Palette: https://coolors.co/palette/003049-d62828-f77f00-fcbf49-eae2b7
var (
	colorBackground = color.RGBA{0, 48, 73, 255}     // ~ Prussian Blue
	colorSnake1     = color.RGBA{252, 191, 73, 255}  // ~ Maximum Yellow Red
	colorSnake2     = color.RGBA{247, 127, 0, 255}   // ~ Orange
	colorFood       = color.RGBA{214, 40, 40, 255}   // ~ Maximum Red
	colorDebug      = color.RGBA{234, 226, 183, 255} // ~ Lemon Meringue
	colorScore      = color.RGBA{247, 127, 0, 255}   // ~ Orange
)

// Game implements ebiten.gameScene interface.
type Game struct {
	curScene scene
}

func NewGame() *Game {
	return &Game{
		curScene: newTitleScene(),
	}
}

// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {

	sceneEnd := g.curScene.update()
	if sceneEnd {
		switch g.curScene.(type) {
		case *titleScene:
			g.curScene = newGameScene()
		}
	}

	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	g.curScene.draw(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
