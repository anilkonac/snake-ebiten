package main

import (
	g "github.com/anilkonac/snake-ebiten/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(g.ScreenSize())
	ebiten.SetWindowTitle("Ssnake")
	ebiten.RunGame(g.NewGame())
}
