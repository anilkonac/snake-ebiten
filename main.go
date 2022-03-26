package main

import (
	g "github.com/anilkonac/snake-ebiten/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(g.ScreenWidth, g.ScreenHeight)
	ebiten.SetWindowTitle("Ssnake")
	ebiten.RunGame(g.NewGame())
}

//go:generate file2byteslice -input game/shaders/basic.go -output game/shaders/basic_go.go -package shaders -var Basic
