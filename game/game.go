package game

import (
	"errors"
	"math/rand"

	"github.com/anilkonac/snake-ebiten/game/object/snake"
	"github.com/anilkonac/snake-ebiten/game/param"
	"github.com/anilkonac/snake-ebiten/game/shader"
	t "github.com/anilkonac/snake-ebiten/game/tool"
	"github.com/hajimehoshi/ebiten/v2"
)

var leadSnake *snake.Snake

const maxTicks = 500

var NumTicks uint16

func init() {
	// rand.Seed(time.Now().UnixNano())
	rand.Seed(1)
}

type scene interface {
	update() bool // Return true if the scene is finished
	draw(*ebiten.Image)
}

// Game implements ebiten.Game interface.
type Game struct {
	curScene scene
}

func NewGame() *Game {
	param.ShaderRound = t.NewShader(shader.Round)

	return &Game{
		curScene: newTitleScene(),
	}
}

// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {

	if NumTicks >= maxTicks {
		return errors.New("bitti")
	}

	if g.curScene.update() {
		switch g.curScene.(type) {
		case *titleScene:
			g.curScene = newGameScene(leadSnake)
		}
	}

	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	g.curScene.draw(screen)
	NumTicks++
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return param.ScreenWidth, param.ScreenHeight
}

func ScreenSize() (int, int) {
	return param.ScreenWidth, param.ScreenHeight
}
