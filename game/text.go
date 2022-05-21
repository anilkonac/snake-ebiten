package game

import (
	"fmt"
	"image"

	"github.com/anilkonac/snake-ebiten/game/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	dpi             = 72
	fontSizeScore   = 32
	fontSizeDebug   = 20
	scoreTextShiftX = 10
	scoreTextShiftY = 8
	fpsTextShiftX   = 0
	fpsTextShiftY   = 3
)

var (
	fontFaceScore  font.Face
	fontFaceDebug  font.Face
	boundScoreText image.Rectangle
	boundFPSText   image.Rectangle
)

func init() {
	tt, err := opentype.Parse(fonts.Rounded)
	if err != nil {
		panic(err)
	}

	fontFaceScore, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeScore,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	initScoreAnim()
	boundScoreText = text.BoundString(fontFaceScore, "Score: 55555")

	tt, err = opentype.Parse(fonts.Debug)
	if err != nil {
		panic(err)
	}

	fontFaceDebug, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeDebug,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	boundFPSText = text.BoundString(fontFaceDebug, "TPS: 60.0\tFPS: 165.0")
}

func drawFPS(screen *ebiten.Image) {
	if printFPS {
		msg := fmt.Sprintf("TPS: %.1f\tFPS: %.1f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
		text.Draw(screen, msg, fontFaceDebug, ScreenWidth-boundFPSText.Size().X-fpsTextShiftX, -boundFPSText.Min.Y+fpsTextShiftY, colorDebug)
	}
}
