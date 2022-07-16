package resource

import (
	"embed"
)

const (
	PathMusic       = "audio/ByeByeBrain_128kbps.ogg"
	PathFontRounded = "fonts/Rounded.ttf"
	PathFontDebug   = "fonts/VT323-Regular.ttf"
)

//go:embed audio fonts
var FS embed.FS
