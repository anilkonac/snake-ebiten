/*
snake-ebiten
Copyright (C) 2022 Anıl Konaç

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	foodScore      = 100
	foodScoreMsg   = "+100"
	decrementAlpha = 4
	decrementY     = 1
)

type scoreAnim struct {
	x, y  int
	color color.RGBA
}

func newScoreAnim(x, y int) *scoreAnim {
	return &scoreAnim{
		x:     x,
		y:     y,
		color: colorScore,
	}
}

func (s *scoreAnim) update() (finished bool) {
	s.y -= decrementY
	if int(s.color.A)-decrementAlpha <= 0 {
		finished = true
	}
	s.color.A -= decrementAlpha
	return
}

func (s *scoreAnim) draw(dst *ebiten.Image) {
	text.Draw(dst, foodScoreMsg, fontScore, s.x-boundScoreAnim.Max.X/2.0, -boundScoreAnim.Min.Y/2.0+s.y, s.color)
}
