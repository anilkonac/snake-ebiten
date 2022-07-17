package resource

import (
	"embed"
)

const (
	PathMusic        = "audio/ByeByeBrain_128kbps.ogg"
	PathSoundEating1 = "audio/eating1.wav"
	PathSoundEating2 = "audio/eating2.wav"
	PathSoundHit     = "audio/hit.wav"
	PathFontRounded  = "fonts/Rounded.ttf"
	PathFontDebug    = "fonts/VT323-Regular.ttf"
)

//go:embed audio/*.wav audio/*.ogg fonts/*.ttf
var FS embed.FS
