/*
Copyright (C) 2022 Anıl Konaç

This file is part of snake-ebiten.

snake-ebiten is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

snake-ebiten is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with snake-ebiten. If not, see <https://www.gnu.org/licenses/>.
*/

package core

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var imagePixel = ebiten.NewImage(1, 1)

func init() {
	imagePixel.Fill(color.White)
}

// TeleCompImage is a TeleComp with drawing options for each rectangle to be used in the DrawImage method.
type TeleCompImage struct {
	TeleComp
	DrawOpts [4]ebiten.DrawImageOptions
}

func (t *TeleCompImage) Update(pureRect *RectF32) {
	t.TeleComp.Update(pureRect)

	for iRect := uint8(0); iRect < t.NumRects; iRect++ {
		rect := &t.Rects[iRect]
		drawOpt := &t.DrawOpts[iRect]

		pos64 := rect.Pos.To64()
		size64 := rect.Size.To64()

		drawOpt.GeoM.Reset()
		drawOpt.GeoM.Scale(size64.X, size64.Y)
		drawOpt.GeoM.Translate(pos64.X, pos64.Y)
	}
}

func (t *TeleCompImage) SetColor(clr color.Color) {
	for iDrawOpt := 0; iDrawOpt < 4; iDrawOpt++ {
		drawOpt := &t.DrawOpts[iDrawOpt]
		drawOpt.ColorM.Reset()
		drawOpt.ColorM.ScaleWithColor(clr)
	}
}

func (t *TeleCompImage) Draw(dst *ebiten.Image) {
	for iRect := uint8(0); iRect < t.NumRects; iRect++ {
		drawOpt := &t.DrawOpts[iRect]
		dst.DrawImage(imagePixel, drawOpt)
	}
}
