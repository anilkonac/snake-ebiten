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
	decrementAlpha = 8
	scoreAnimSpeed = 20
)

var (
	scoreAnimShiftX      float32
	scoreAnimShiftY      float32
	scoreAnimBound       image.Rectangle
	scoreAnimBoundSize   image.Point
	foodScoreMsg         = strconv.Itoa(foodScore)
	scoreAnimImage       *ebiten.Image
	drawOptionsScoreAnim ebiten.DrawTrianglesShaderOptions
)

type scoreAnim struct {
	x, y      float32
	alpha     uint8
	direction directionT
	rects     []rectF32
}

func initScoreAnim() {
	scoreAnimBound = text.BoundString(fontScore, foodScoreMsg)
	scoreAnimBoundSize = scoreAnimBound.Size()

	scoreAnimShiftX = halfSnakeWidth + float32(scoreAnimBoundSize.X)/2.0
	scoreAnimShiftY = halfSnakeWidth + float32(scoreAnimBoundSize.Y)/2.0

	// Prepare score animation text image.
	scoreAnimImage = ebiten.NewImage(scoreAnimBoundSize.X, scoreAnimBoundSize.Y)
	scoreAnimImage.Fill(color.Black)
	text.Draw(scoreAnimImage, foodScoreMsg, fontScore,
		-scoreAnimBound.Min.X, -scoreAnimBound.Min.Y,
		color.RGBA{255, 0, 0, 255})

	// Prepare draw options
	drawOptionsScoreAnim = ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]interface{}{
			"RawAlpha": float32(1.0),
		},
		Images: [4]*ebiten.Image{scoreAnimImage, nil, nil, nil},
	}
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

	newAnim := &scoreAnim{
		x:         x,
		y:         y,
		alpha:     colorScore.A,
		direction: dir,
	}

	newAnim.createRects()

	return newAnim
}

func (s *scoreAnim) createRects() {
	// Create a rectangle to be split
	pureRect := rectF32{
		x:      s.x - float32(scoreAnimBoundSize.X)/2.0,
		y:      s.y - float32(scoreAnimBoundSize.Y)/2.0,
		width:  float32(scoreAnimBoundSize.X),
		height: float32(scoreAnimBoundSize.Y),
	}
	// Init/Remove rects
	s.rects = make([]rectF32, 0, 4) // Remove rects

	// Split this rectangle if it is on a screen edge.
	pureRect.split(&s.rects)
}

func (s *scoreAnim) update() (finished bool) {
	// Move animation coordinate
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

	// Update rectangles of this anim
	s.createRects()

	// Update alpha
	if int(s.alpha)-decrementAlpha <= 0 {
		finished = true
	}
	s.alpha -= decrementAlpha
	drawOptionsScoreAnim.Uniforms["RawAlpha"] = float32(s.alpha)

	return
}

// Implement drawable interface
// ----------------------------
func (s *scoreAnim) drawEnabled() bool {
	return true
}

func (s *scoreAnim) drawableRects() []rectF32 {
	return s.rects
}

func (s *scoreAnim) Color() color.Color {
	return colorScore
}

func (s *scoreAnim) drawOptions() *ebiten.DrawTrianglesShaderOptions {
	return &drawOptionsScoreAnim
}

func (s *scoreAnim) shader() *ebiten.Shader {
	return shaderScore
}

func (s *scoreAnim) drawDebugInfo(dst *ebiten.Image) {

}
