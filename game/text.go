package game

import (
	"github.com/anilkonac/snake-ebiten/game/resources/fonts"
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
	fontScore font.Face
	fontDebug font.Face
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
	initScoreAnim()

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
