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
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const mouthRadius float32 = halfSnakeWidth * 0.75

type unit struct {
	headCenter     vec64
	length         float64
	direction      directionT
	rectsCollision []rectF32
	rectsDrawable  []rectF32
	color          *color.RGBA
	next           *unit
	prev           *unit
	drawOpts       ebiten.DrawTrianglesShaderOptions
}

func newUnit(headCenter vec64, length float64, direction directionT, color *color.RGBA) *unit {
	newUnit := &unit{
		headCenter: headCenter,
		length:     length,
		direction:  direction,
		color:      color,
		drawOpts: ebiten.DrawTrianglesShaderOptions{
			Uniforms: map[string]interface{}{
				"Radius":      float32(halfSnakeWidth),
				"RadiusMouth": mouthRadius,
			},
		},
	}
	newUnit.update(eatingAnimStartDistance)

	return newUnit
}

func (u *unit) createRects() {
	// Create rectangles for drawing and collision. They are going to split.
	var rectDraw, rectColl *rectF32

	rectColl = u.createRectColl()
	rectDraw = u.createRectDraw(rectColl)

	// Remove old rectangles
	u.rectsCollision = make([]rectF32, 0, 4)
	if u.next != nil {
		u.rectsDrawable = make([]rectF32, 0, 4)
	}

	// Create split rectangles on screen edges.
	rectColl.split(&u.rectsCollision)
	if u.next == nil {
		u.rectsDrawable = u.rectsCollision
		return
	}
	rectDraw.split(&u.rectsDrawable)
}

func (u *unit) createRectColl() (rectColl *rectF32) {
	length32 := float32(math.Floor(u.length))
	flCenter := u.headCenter.floor().to32()

	switch u.direction {
	case directionRight:
		rectColl = newRect(vec32{flCenter.x - length32 + halfSnakeWidth, flCenter.y - halfSnakeWidth}, vec32{length32, snakeWidth})
	case directionLeft:
		rectColl = newRect(vec32{flCenter.x - halfSnakeWidth, flCenter.y - halfSnakeWidth}, vec32{length32, snakeWidth})
	case directionUp:
		rectColl = newRect(vec32{flCenter.x - halfSnakeWidth, flCenter.y - halfSnakeWidth}, vec32{snakeWidth, length32})
	case directionDown:
		rectColl = newRect(vec32{flCenter.x - halfSnakeWidth, flCenter.y - length32 + halfSnakeWidth}, vec32{snakeWidth, length32})
	default:
		panic("Wrong unit direction!!")
	}

	return
}

func (u *unit) createRectDraw(rectColl *rectF32) (rectDraw *rectF32) {
	if u.next == nil {
		rectDraw = rectColl
		return
	}

	switch u.direction {
	case directionRight:
		rectDraw = newRect(vec32{rectColl.pos.x - snakeWidth, rectColl.pos.y}, vec32{rectColl.size.x + snakeWidth, rectColl.size.y})
	case directionLeft:
		rectDraw = newRect(vec32{rectColl.pos.x, rectColl.pos.y}, vec32{rectColl.size.x + snakeWidth, rectColl.size.y})
	case directionUp:
		rectDraw = newRect(vec32{rectColl.pos.x, rectColl.pos.y}, vec32{rectColl.size.x, rectColl.size.y + snakeWidth})
	case directionDown:
		rectDraw = newRect(vec32{rectColl.pos.x, rectColl.pos.y - snakeWidth}, vec32{rectColl.size.x, rectColl.size.y + snakeWidth})
	default:
		panic("Wrong unit direction!!")
	}

	return
}

func (u *unit) update(distToFood float32) {
	u.createRects()                 // Update rectangles of this unit
	u.updateDrawOptions(distToFood) // Update draw options
}

func (u *unit) updateDrawOptions(distToFood float32) {
	// Distance to food
	uniDistFood := distToFood / eatingAnimStartDistance

	// Specify Size uniform variable
	var drawWidth, drawHeight float32
	flooredLength := float32(math.Floor(u.length))
	if u.next != nil {
		flooredLength += snakeWidth
	}
	if u.direction.isVertical() {
		drawWidth, drawHeight = snakeWidth, flooredLength
	} else {
		drawWidth, drawHeight = flooredLength, snakeWidth
	}

	// Update the options
	u.drawOpts.Uniforms["Direction"] = float32(u.direction)
	u.drawOpts.Uniforms["Size"] = []float32{drawWidth, drawHeight}
	u.drawOpts.Uniforms["DistToFood"] = uniDistFood
}

func (u *unit) moveUp(dist float64) {
	u.headCenter.y -= dist

	// teleport if head center is offscreen.
	if u.headCenter.y < 0 {
		u.headCenter.y += ScreenHeight
	}
}

func (u *unit) moveDown(dist float64) {
	u.headCenter.y += dist

	// teleport if head center is offscreen.
	if u.headCenter.y > ScreenHeight {
		u.headCenter.y -= ScreenHeight
	}
}

func (u *unit) moveRight(dist float64) {
	u.headCenter.x += dist

	// teleport if head center is offscreen.
	if u.headCenter.x > ScreenWidth {
		u.headCenter.x -= ScreenWidth
	}
}

func (u *unit) moveLeft(dist float64) {
	u.headCenter.x -= dist

	// teleport if head center is offscreen.
	if u.headCenter.x < 0 {
		u.headCenter.x += ScreenWidth
	}
}

func (u *unit) markHeadCenters(dst *ebiten.Image) {
	markPoint(dst, &u.headCenter, 4, colorFood)

	var offset float64 = 0
	if u.next == nil {
		offset = snakeWidth
	}

	backCenter := u.headCenter
	switch u.direction {
	case directionUp:
		backCenter.y = u.headCenter.y + u.length - offset
	case directionDown:
		backCenter.y = u.headCenter.y - u.length + offset
	case directionRight:
		backCenter.x = u.headCenter.x - u.length + offset
	case directionLeft:
		backCenter.x = u.headCenter.x + u.length - offset
	}
	// mark head center at the other side
	markPoint(dst, &backCenter, 4, colorFood)
}

// Implement collidable interface
// ------------------------------
func (u *unit) collEnabled() bool {
	return true
}

func (u *unit) collisionRects() []rectF32 {
	return u.rectsCollision
}

// Implement drawable interface
// ----------------------------
func (u *unit) drawEnabled() bool {
	return true
}

func (u *unit) drawableRects() []rectF32 {
	return u.rectsDrawable
}

func (u *unit) Color() *color.RGBA {
	return u.color
}

func (u *unit) drawOptions() *ebiten.DrawTrianglesShaderOptions {
	return &u.drawOpts
}

func (u *unit) shader() *ebiten.Shader {
	if u.prev == nil {
		return shaderSnakeHead
	}
	return shaderRound
}

func (u *unit) drawDebugInfo(dst *ebiten.Image) {
	u.markHeadCenters(dst)
	for iRect := range u.rectsDrawable {
		rect := u.rectsDrawable[iRect]
		rect.drawOuterRect(dst, colorFood)
	}
}
