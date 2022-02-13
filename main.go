package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ssnake")

	// Start game loop
	game := newGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
