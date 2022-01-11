package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game implements ebiten.Game interface.
type Game struct {
	mouseX int
	mouseY int
}

// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	g.mouseX, g.mouseY = ebiten.CursorPosition()
	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Mouse X: %v, Mouse Y: %v", g.mouseX, g.mouseY))
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// s := ebiten.DeviceScaleFactor()
	// return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
	return screenWidth, screenHeight
}
