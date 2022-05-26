package game

import (
	"github.com/anilkonac/snake-ebiten/game/params"
	"github.com/anilkonac/snake-ebiten/game/shaders"
	t "github.com/anilkonac/snake-ebiten/game/tools"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game implements ebiten.Game interface.
type Game struct {
	curScene scene
}

func init() {
	params.ShaderRound = t.NewShader(shaders.Round)
}

func NewGame() *Game {
	return &Game{
		curScene: newTitleScene(),
	}
}

// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	if g.curScene.update() {
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
	return params.ScreenWidth, params.ScreenHeight
}

func ScreenSize() (int, int) {
	return params.ScreenWidth, params.ScreenHeight
}
