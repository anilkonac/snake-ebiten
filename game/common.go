package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func drawFPS(screen *ebiten.Image) {
	if printFPS {
		msg := fmt.Sprintf("TPS: %.1f\tFPS: %.1f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
		text.Draw(screen, msg, fontFaceDebug, ScreenWidth-boundFPSText.Size().X-fpsTextShiftX, -boundFPSText.Min.Y+fpsTextShiftY, colorDebug)
	}
}
