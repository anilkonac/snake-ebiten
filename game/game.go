package game

import "github.com/hajimehoshi/ebiten/v2"

const (
	ScreenWidth  = 960
	ScreenHeight = 720
	deltaTime    = 1.0 / 60.0
)

// Game implements ebiten.gameScene interface.
type Game struct {
	curScene scene
}

func NewGame() *Game {
	return &Game{
		curScene: newTitleScreen(),
	}
}

// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	g.curScene.update()

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
