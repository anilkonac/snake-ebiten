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

package object

import (
	"image/color"
	"strconv"

	c "github.com/anilkonac/snake-ebiten/game/core"
	s "github.com/anilkonac/snake-ebiten/game/object/snake"
	"github.com/anilkonac/snake-ebiten/game/param"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	decrementAlpha   = 8.0
	scoreAnimSpeed   = 25
	scoreAnimPadding = 8
)

var (
	scoreAnimShiftY    float32
	scoreAnimBoundSize c.Vec32
	scoreAnimImage     *ebiten.Image
)

type ScoreAnim struct {
	c.TeleCompTriang
	pos       c.Vec32
	alpha     uint8
	direction s.DirectionT
	drawOpts  ebiten.DrawTrianglesOptions
}

func InitScoreAnim() {
	// Init animation text bound variables
	foodScoreMsg := strconv.Itoa(param.FoodScore)
	scoreAnimBound := text.BoundString(param.FontFaceScore, foodScoreMsg)
	scoreAnimBoundSizeI := scoreAnimBound.Size()
	scoreAnimBoundSize.X = float32(scoreAnimBoundSizeI.X)
	scoreAnimBoundSize.Y = float32(scoreAnimBoundSizeI.Y)
	scoreAnimShiftY = param.RadiusSnake + scoreAnimBoundSize.Y/2.0 + scoreAnimPadding

	// Prepare score animation text image.
	scoreAnimImage = ebiten.NewImage(scoreAnimBoundSizeI.X, scoreAnimBoundSizeI.Y)
	text.Draw(scoreAnimImage, foodScoreMsg, param.FontFaceScore,
		-scoreAnimBound.Min.X, -scoreAnimBound.Min.Y,
		color.White)
}

func NewScoreAnim(pos c.Vec32) *ScoreAnim {
	newAnim := &ScoreAnim{
		pos: c.Vec32{
			X: pos.X,
			Y: pos.Y - scoreAnimShiftY,
		},
		alpha:     param.ColorScore.A,
		direction: s.DirectionUp,
	}
	newAnim.SetColor(&param.ColorScore)

	newAnim.createRects()

	return newAnim
}

func (s *ScoreAnim) createRects() {
	// Create a rectangle to be split
	pureRect := c.RectF32{
		Pos: c.Vec32{
			X: s.pos.X - scoreAnimBoundSize.X/2.0,
			Y: s.pos.Y - scoreAnimBoundSize.Y/2.0,
		},
		Size: c.Vec32{X: scoreAnimBoundSize.X, Y: scoreAnimBoundSize.Y},
	}
	// Split this rectangle if it is on a screen edge.
	s.TeleCompTriang.Update(&pureRect)
}

// Returns true when the animation is finished
func (s *ScoreAnim) Update() bool {
	// Move animation
	s.pos.Y -= scoreAnimSpeed * param.DeltaTime

	// Decrease alpha
	if s.alpha < decrementAlpha {
		return true
	}

	s.alpha -= decrementAlpha
	color := color.RGBA{param.ColorScore.R, param.ColorScore.G, param.ColorScore.B, s.alpha}
	s.SetColor(&color)

	// Update rectangles of this anim
	s.createRects()

	return false
}

func (s *ScoreAnim) Draw(dst *ebiten.Image) {
	vertices, indices := s.Triangles()
	dst.DrawTriangles(vertices, indices, scoreAnimImage, &s.drawOpts)
}
