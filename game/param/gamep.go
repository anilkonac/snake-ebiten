package param

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

const (
	ScreenWidth      = 960
	ScreenHeight     = 720
	HalfScreenWidth  = ScreenWidth / 2.0
	HalfScreenHeight = ScreenHeight / 2.0
	DeltaTime        = 1.0 / 60.0
)

// Food parameters
const (
	FoodScore    = 100
	FoodLength   = 16
	RadiusFood   = FoodLength / 2.0
	RadiusEating = RadiusMouth + RadiusFood
)

// Colors to be used in the drawing.
// Palette: https://coolors.co/palette/003049-d62828-f77f00-fcbf49-eae2b7
var (
	ColorBackground = color.RGBA{0, 48, 73, 255}     // ~ Prussian Blue
	ColorSnake1     = color.RGBA{252, 191, 73, 255}  // ~ Maximum Yellow Red
	ColorSnake2     = color.RGBA{247, 127, 0, 255}   // ~ Orange
	ColorFood       = color.RGBA{214, 40, 40, 255}   // ~ Maximum Red
	ColorDebug      = color.RGBA{234, 226, 183, 255} // ~ Lemon Meringue
	ColorScore      = color.RGBA{247, 127, 0, 255}   // ~ Orange
)

var (
	TeleportEnabled = true
	PrintFPS        = true
	DebugUnits      = false // Draw consecutive units with different colors
	ShaderRound     *ebiten.Shader
	FontFaceScore   font.Face
)
