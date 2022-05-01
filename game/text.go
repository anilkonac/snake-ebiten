package game

import (
	"image"

	"github.com/anilkonac/snake-ebiten/game/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	dpi             = 84
	fontSizeScore   = 24
	fontSizeDebug   = 16
	scoreTextShiftX = 10
	scoreTextShiftY = 8
	fpsTextShiftX   = 3
	fpsTextShiftY   = 2
)

var (
	fontScore      font.Face
	fontDebug      font.Face
	boundScoreAnim image.Rectangle
)

func init() {
	tt, err := opentype.Parse(fonts.Rounded)
	if err != nil {
		panic(err)
	}

	fontScore, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeScore,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	boundScoreAnim = text.BoundString(fontScore, foodScoreMsg)

	tt, err = opentype.Parse(fonts.Debug)
	if err != nil {
		panic(err)
	}

	fontDebug, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeDebug,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
}
