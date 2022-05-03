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
	"image"
	"image/color"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	foodScore      = 100
	decrementAlpha = 6
	scoreAnimSpeed = 20
)

var (
	scoreAnimShiftX    float32
	scoreAnimShiftY    float32
	scoreAnimBound     image.Rectangle
	scoreAnimBoundSize image.Point
	foodScoreMsg       = strconv.Itoa(foodScore)
)

type scoreAnim struct {
	x, y      float32
	alpha     uint8
	direction directionT
}

func initScoreAnim() {
	scoreAnimBound = text.BoundString(fontScore, foodScoreMsg)
	scoreAnimBoundSize = scoreAnimBound.Size()
	scoreAnimShiftX = halfSnakeWidth + float32(scoreAnimBoundSize.X)/2.0
	scoreAnimShiftY = halfSnakeWidth + float32(scoreAnimBoundSize.Y)/2.0
}

func newScoreAnim(x, y float32, verticalDir bool) *scoreAnim {
	// Determine the direction of the new animation
	var dir directionT
	if tossUp := rand.Intn(2); verticalDir {
		dir = directionT(tossUp)
		// Shift the animation near the snake
		if dir == directionUp {
			y -= scoreAnimShiftY
		} else {
			y += scoreAnimShiftY
		}
	} else {
		dir = 2 + directionT(tossUp)
		// Shift the animation near the snake
		if dir == directionLeft {
			x -= scoreAnimShiftX
		} else {
			x += scoreAnimShiftX
		}
	}
	return &scoreAnim{x, y, colorScore.A, dir}
}

func (s *scoreAnim) update() (finished bool) {
	switch s.direction {
	case directionUp:
		s.y -= scoreAnimSpeed * deltaTime
	case directionDown:
		s.y += scoreAnimSpeed * deltaTime
	case directionLeft:
		s.x -= scoreAnimSpeed * deltaTime
	case directionRight:
		s.x += scoreAnimSpeed * deltaTime
	}

	if int(s.alpha)-decrementAlpha <= 0 {
		finished = true
	}
	s.alpha -= decrementAlpha
	return
}

func (s *scoreAnim) draw(dst *ebiten.Image) {
	text.Draw(dst, foodScoreMsg, fontScore,
		int(s.x)-(scoreAnimBound.Max.X/2.0), int(s.y)-(scoreAnimBound.Min.Y/2.0),
		color.RGBA{colorScore.R, colorScore.G, colorScore.B, s.alpha})
}
