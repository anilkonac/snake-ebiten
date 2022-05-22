package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type inputHandler struct {
	keys []ebiten.Key
}

func (i *inputHandler) getPressedKeys() {
	i.keys = inpututil.AppendPressedKeys(i.keys[:0])
}
