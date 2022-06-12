package game

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/anilkonac/snake-ebiten/game/object/snake"
	"github.com/anilkonac/snake-ebiten/game/param"
	"github.com/anilkonac/snake-ebiten/game/shader"
	t "github.com/anilkonac/snake-ebiten/game/tool"
	"github.com/hajimehoshi/ebiten/v2"
)

var leadSnake *snake.Snake

const maxTicks = 500

var (
	measureTime    bool
	NumTicks       uint16
	FPSSum, TPSSum float64
	TPSMin, FPSMin float64 = math.MaxFloat64, math.MaxFloat64
	TPSMax, FPSMax float64
)

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

	go func() {
		time.Sleep(time.Second * 3)
		measureTime = true
	}()

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

	if !measureTime {
		return
	}

	NumTicks++
	tps, fps := ebiten.CurrentTPS(), ebiten.CurrentFPS()
	TPSSum += tps
	FPSSum += fps
	TPSMin = math.Min(TPSMin, tps)
	TPSMax = math.Max(TPSMax, tps)
	FPSMin = math.Min(FPSMin, fps)
	FPSMax = math.Max(FPSMax, fps)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return param.ScreenWidth, param.ScreenHeight
}

func ScreenSize() (int, int) {
	return param.ScreenWidth, param.ScreenHeight
}
