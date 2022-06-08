package game

import (
	"fmt"
	"image"

	c "github.com/anilkonac/snake-ebiten/game/core"
	"github.com/anilkonac/snake-ebiten/game/object"
	"github.com/anilkonac/snake-ebiten/game/param"
	"github.com/anilkonac/snake-ebiten/resources/fonts"
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
	fontFaceDebug      font.Face
	fontFaceTitle      font.Face
	boundTextScore     image.Rectangle
	boundTextFPS       image.Rectangle
	boundTextTitle     image.Rectangle
	boundTextKeyPrompt image.Rectangle
)

func init() {
	tt, err := opentype.Parse(fonts.Rounded)
	c.Panic(err)

	param.FontFaceScore, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeScore,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	c.Panic(err)

	fontFaceTitle, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeTitle,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	c.Panic(err)

	tt, err = opentype.Parse(fonts.Debug)
	c.Panic(err)

	fontFaceDebug, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeDebug,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	c.Panic(err)

	boundTextScore = text.BoundString(param.FontFaceScore, "Score: 55555")
	boundTextTitle = text.BoundString(fontFaceTitle, textTitle)
	boundTextKeyPrompt = text.BoundString(param.FontFaceScore, textPressToPlay)
	boundTextFPS = text.BoundString(fontFaceDebug, "TPS: 60.0\tFPS: 165.0")

	object.InitScoreAnim()
}

func drawFPS(screen *ebiten.Image) {
	if param.PrintFPS {
		msg := fmt.Sprintf("TPS: %.1f\tFPS: %.1f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
		text.Draw(screen, msg, fontFaceDebug, param.ScreenWidth-boundTextFPS.Size().X-fpsTextShiftX, -boundTextFPS.Min.Y+fpsTextShiftY, param.ColorDebug)
	}
}
