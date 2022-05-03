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
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	foodScore      = 100
	decrementAlpha = 4
	scoreAnimSpeed = 1
	scoreAnimShift = snakeWidth
)

var foodScoreMsg = strconv.Itoa(foodScore)

type scoreAnim struct {
	x, y      int
	alpha     uint8
	direction directionT
}

func newScoreAnim(x, y int, verticalDir bool) *scoreAnim {
	// Determine the direction of the new animation
	var dir directionT
	if tossUp := rand.Intn(2); verticalDir {
		dir = directionT(tossUp)
		// Shift the animation near the snake
		if dir == directionUp {
			y -= scoreAnimShift
		} else {
			y += scoreAnimShift
		}
	} else {
		dir = 2 + directionT(tossUp)
		// Shift the animation near the snake
		if dir == directionLeft {
			x -= scoreAnimShift
		} else {
			x += scoreAnimShift
		}
	}
	return &scoreAnim{x, y, colorScore.A, dir}
}

func (s *scoreAnim) update() (finished bool) {
	switch s.direction {
	case directionUp:
		s.y -= scoreAnimSpeed
	case directionDown:
		s.y += scoreAnimSpeed
	case directionLeft:
		s.x -= scoreAnimSpeed
	case directionRight:
		s.x += scoreAnimSpeed
	}

	if int(s.alpha)-decrementAlpha <= 0 {
		finished = true
	}
	s.alpha -= decrementAlpha
	return
}

func (s *scoreAnim) draw(dst *ebiten.Image) {
	text.Draw(dst, foodScoreMsg, fontScore, s.x-(boundScoreAnim.Max.X/2.0), -(boundScoreAnim.Min.Y/2.0)+s.y, color.RGBA{colorScore.R, colorScore.G, colorScore.B, s.alpha})
}
