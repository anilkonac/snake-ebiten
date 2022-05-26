package game

import (
	"fmt"
	"image"

	"github.com/anilkonac/snake-ebiten/game/params"
	"github.com/anilkonac/snake-ebiten/game/resources/fonts"
	t "github.com/anilkonac/snake-ebiten/game/tools"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	dpi             = 72
	fontSizeScore   = 32
	fontSizeDebug   = 20
	fontSizeTitle   = 128
	scoreTextShiftX = 10
	scoreTextShiftY = 8
	fpsTextShiftX   = 0
	fpsTextShiftY   = 3
)

var (
	fontFaceScore      font.Face
	fontFaceDebug      font.Face
	fontFaceTitle      font.Face
	boundTextScore     image.Rectangle
	boundTextFPS       image.Rectangle
	boundTextTitle     image.Rectangle
	boundTextKeyPrompt image.Rectangle
)

func init() {
	tt, err := opentype.Parse(fonts.Rounded)
	t.Panic(err)

	fontFaceScore, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeScore,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	t.Panic(err)

	fontFaceTitle, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeTitle,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	t.Panic(err)

	tt, err = opentype.Parse(fonts.Debug)
	t.Panic(err)

	fontFaceDebug, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeDebug,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	t.Panic(err)

	boundTextScore = text.BoundString(fontFaceScore, "Score: 55555")
	boundTextTitle = text.BoundString(fontFaceTitle, textTitle)
	boundTextKeyPrompt = text.BoundString(fontFaceScore, textPressToPlay)
	boundTextFPS = text.BoundString(fontFaceDebug, "TPS: 60.0\tFPS: 165.0")

	initScoreAnim()
}

func drawFPS(screen *ebiten.Image) {
	if printFPS {
		msg := fmt.Sprintf("TPS: %.1f\tFPS: %.1f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
		text.Draw(screen, msg, fontFaceDebug, params.ScreenWidth-boundTextFPS.Size().X-fpsTextShiftX, -boundTextFPS.Min.Y+fpsTextShiftY, params.ColorDebug)
	}
}
