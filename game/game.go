package game

import (
	"github.com/anilkonac/snake-ebiten/game/object/snake"
	"github.com/anilkonac/snake-ebiten/game/param"
	"github.com/anilkonac/snake-ebiten/game/shader"
	t "github.com/anilkonac/snake-ebiten/game/tool"
	"github.com/hajimehoshi/ebiten/v2"
)

type scene interface {
	update() bool // Return true if the scene is finished
	draw(*ebiten.Image)
}

// Game implements ebiten.Game interface.
type Game struct {
	curScene    scene
	playerSnake *snake.Snake
}

func NewGame() *Game {
	param.ShaderRound = t.NewShader(shader.Round)
	playerSnake := snake.NewSnakeRandDir(
		t.Vec64{X: snakeHeadCenterX, Y: snakeHeadCenterY},
		param.SnakeLength, param.SnakeSpeedInitial, &param.ColorSnake1,
	)

	return &Game{
		curScene:    newTitleScene(playerSnake),
		playerSnake: playerSnake,
	}
}

// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	if g.curScene.update() {
		switch g.curScene.(type) {
		case *titleScene:
			g.curScene = newGameScene(g.playerSnake)
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
	return param.ScreenWidth, param.ScreenHeight
}

func ScreenSize() (int, int) {
	return param.ScreenWidth, param.ScreenHeight
}
