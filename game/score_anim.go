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
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	foodScore         = 100
	decrementAlpha    = 8.0 / 255.0
	scoreAnimSpeed    = 20
	relativeRandomDir = false
)

var (
	scoreAnimShiftY      float32
	scoreAnimBoundWidth  float32
	scoreAnimBoundHeight float32
	scoreAnimImage       *ebiten.Image
	drawOptionsScoreAnim ebiten.DrawTrianglesShaderOptions
)

type scoreAnim struct {
	x, y      float32
	alpha     float32
	direction directionT
	rects     []rectF32
}

func initScoreAnim() {
	// Init animation text bound variables
	foodScoreMsg := strconv.Itoa(foodScore)
	scoreAnimBound := text.BoundString(fontScore, foodScoreMsg)
	scoreAnimBoundSize := scoreAnimBound.Size()
	scoreAnimBoundWidth = float32(scoreAnimBoundSize.X)
	scoreAnimBoundHeight = float32(scoreAnimBoundSize.Y)
	scoreAnimShiftY = halfSnakeWidth + scoreAnimBoundHeight/2.0

	// Prepare score animation text image.
	scoreAnimImage = ebiten.NewImage(scoreAnimBoundSize.X, scoreAnimBoundSize.Y)
	scoreAnimImage.Fill(color.Black)
	text.Draw(scoreAnimImage, foodScoreMsg, fontScore,
		-scoreAnimBound.Min.X, -scoreAnimBound.Min.Y,
		color.RGBA{255, 0, 0, 255})

	// Prepare draw options
	drawOptionsScoreAnim.Uniforms = map[string]interface{}{
		"Alpha": float32(1.0),
	}
	drawOptionsScoreAnim.Images = [4]*ebiten.Image{scoreAnimImage, nil, nil, nil}
}

func newScoreAnim(x, y float32, verticalDir bool) *scoreAnim {
	newAnim := &scoreAnim{
		x:         x,
		y:         y - scoreAnimShiftY,
		alpha:     float32(colorScore.A) / 255.0,
		direction: directionUp,
	}

	newAnim.createRects()

	return newAnim
}

func (s *scoreAnim) createRects() {
	// Create a rectangle to be split
	pureRect := rectF32{
		x:      s.x - scoreAnimBoundWidth/2.0,
		y:      s.y - scoreAnimBoundHeight/2.0,
		width:  scoreAnimBoundWidth,
		height: scoreAnimBoundHeight,
	}
	// Init/Remove rects
	s.rects = make([]rectF32, 0, 4) // Remove rects

	// Split this rectangle if it is on a screen edge.
	pureRect.split(&s.rects)
}

// Returns true when the animation is finished
func (s *scoreAnim) update() bool {
	// Move animation
	s.y -= scoreAnimSpeed * deltaTime

	// Update rectangles of this anim
	s.createRects()

	// Decrease alpha
	s.alpha -= decrementAlpha
	drawOptionsScoreAnim.Uniforms["Alpha"] = s.alpha

	return (s.alpha <= 0.0)
}

// Implement drawable interface
// ----------------------------
func (s *scoreAnim) drawEnabled() bool {
	return true
}

func (s *scoreAnim) drawableRects() []rectF32 {
	return s.rects
}

func (s *scoreAnim) Color() *color.RGBA {
	return &colorScore
}

func (s *scoreAnim) drawOptions() *ebiten.DrawTrianglesShaderOptions {
	return &drawOptionsScoreAnim
}

func (s *scoreAnim) shader() *ebiten.Shader {
	return shaderScore
}

func (s *scoreAnim) drawDebugInfo(dst *ebiten.Image) {

}
