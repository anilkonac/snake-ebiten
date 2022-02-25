package main

import (
	g "github.com/anilkonac/snake-ebiten/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

func main() {
	// ebiten.SetWindowSize(g.GameWidth, g.GameHeight)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ssnake")
	ebiten.RunGame(g.NewGame())
}
