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

	"github.com/anilkonac/snake-ebiten/game/params"
	"github.com/anilkonac/snake-ebiten/game/shaders"
	t "github.com/anilkonac/snake-ebiten/game/tools"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	decrementAlpha   = 8.0 / 255.0
	scoreAnimSpeed   = 25
	scoreAnimPadding = 8
)

var (
	scoreAnimShiftY      float32
	scoreAnimBoundSize   t.Vec32
	scoreAnimImage       *ebiten.Image
	shaderScore          *ebiten.Shader
	drawOptionsScoreAnim ebiten.DrawTrianglesShaderOptions
)

type scoreAnim struct {
	pos       t.Vec32
	alpha     float32
	direction directionT
	rects     []t.RectF32
}

func initScoreAnim() {
	shaderScore = t.NewShader(shaders.Score)

	// Init animation text bound variables
	foodScoreMsg := strconv.Itoa(params.FoodScore)
	scoreAnimBound := text.BoundString(fontFaceScore, foodScoreMsg)
	scoreAnimBoundSizeI := scoreAnimBound.Size()
	scoreAnimBoundSize.X = float32(scoreAnimBoundSizeI.X)
	scoreAnimBoundSize.Y = float32(scoreAnimBoundSizeI.Y)
	scoreAnimShiftY = params.RadiusSnake + scoreAnimBoundSize.Y/2.0 + scoreAnimPadding

	// Prepare score animation text image.
	scoreAnimImage = ebiten.NewImage(scoreAnimBoundSizeI.X, scoreAnimBoundSizeI.Y)
	scoreAnimImage.Fill(color.Black)
	text.Draw(scoreAnimImage, foodScoreMsg, fontFaceScore,
		-scoreAnimBound.Min.X, -scoreAnimBound.Min.Y,
		color.RGBA{255, 0, 0, 255})

	// Prepare draw options
	drawOptionsScoreAnim.Uniforms = map[string]interface{}{
		"Alpha": float32(1.0),
	}
	drawOptionsScoreAnim.Images = [4]*ebiten.Image{scoreAnimImage, nil, nil, nil}
}

func newScoreAnim(pos t.Vec32) *scoreAnim {
	newAnim := &scoreAnim{
		pos: t.Vec32{
			X: pos.X,
			Y: pos.Y - scoreAnimShiftY,
		},
		alpha:     float32(params.ColorScore.A) / 255.0,
		direction: directionUp,
	}

	newAnim.createRects()

	return newAnim
}

func (s *scoreAnim) createRects() {
	// Create a rectangle to be split
	pureRect := t.RectF32{
		Pos: t.Vec32{
			X: s.pos.X - scoreAnimBoundSize.X/2.0,
			Y: s.pos.Y - scoreAnimBoundSize.Y/2.0,
		},
		Size: t.Vec32{X: scoreAnimBoundSize.X, Y: scoreAnimBoundSize.Y},
	}
	// Init/Remove rects
	s.rects = make([]t.RectF32, 0, 4) // Remove rects

	// Split this rectangle if it is on a screen edge.
	pureRect.Split(&s.rects)
}

// Returns true when the animation is finished
func (s *scoreAnim) update() bool {
	// Move animation
	s.pos.Y -= scoreAnimSpeed * params.DeltaTime

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

func (s *scoreAnim) drawableRects() []t.RectF32 {
	return s.rects
}

func (s *scoreAnim) Color() *color.RGBA {
	return &params.ColorScore
}

func (s *scoreAnim) drawOptions() *ebiten.DrawTrianglesShaderOptions {
	return &drawOptionsScoreAnim
}

func (s *scoreAnim) shader() *ebiten.Shader {
	return shaderScore
}

func (s *scoreAnim) drawDebugInfo(dst *ebiten.Image) {

}
