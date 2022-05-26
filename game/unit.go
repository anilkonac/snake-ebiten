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

	"github.com/anilkonac/snake-ebiten/game/params"
	"github.com/hajimehoshi/ebiten/v2"
)

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
				"Radius":      float32(params.RadiusSnake),
				"RadiusMouth": float32(params.RadiusMouth),
			},
		},
	}
	newUnit.update(params.EatingAnimStartDistance)

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
		rectColl = newRect(vec32{flCenter.x - length32 + params.RadiusSnake, flCenter.y - params.RadiusSnake}, vec32{length32, params.SnakeWidth})
	case directionLeft:
		rectColl = newRect(vec32{flCenter.x - params.RadiusSnake, flCenter.y - params.RadiusSnake}, vec32{length32, params.SnakeWidth})
	case directionUp:
		rectColl = newRect(vec32{flCenter.x - params.RadiusSnake, flCenter.y - params.RadiusSnake}, vec32{params.SnakeWidth, length32})
	case directionDown:
		rectColl = newRect(vec32{flCenter.x - params.RadiusSnake, flCenter.y - length32 + params.RadiusSnake}, vec32{params.SnakeWidth, length32})
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
		rectDraw = newRect(vec32{rectColl.pos.x - params.SnakeWidth, rectColl.pos.y}, vec32{rectColl.size.x + params.SnakeWidth, rectColl.size.y})
	case directionLeft:
		rectDraw = newRect(vec32{rectColl.pos.x, rectColl.pos.y}, vec32{rectColl.size.x + params.SnakeWidth, rectColl.size.y})
	case directionUp:
		rectDraw = newRect(vec32{rectColl.pos.x, rectColl.pos.y}, vec32{rectColl.size.x, rectColl.size.y + params.SnakeWidth})
	case directionDown:
		rectDraw = newRect(vec32{rectColl.pos.x, rectColl.pos.y - params.SnakeWidth}, vec32{rectColl.size.x, rectColl.size.y + params.SnakeWidth})
	default:
		panic("Wrong unit direction!!")
	}

	return
}

func (u *unit) update(distToFood float32) {
	u.createRects() // Update rectangles of this unit
	u.updateDrawOptions(distToFood)
}

func (u *unit) updateDrawOptions(distToFood float32) {
	// Distance to food
	proxToFood := 1.0 - distToFood/params.EatingAnimStartDistance

	// Specify Size uniform variable
	var drawWidth, drawHeight float32
	flooredLength := float32(math.Floor(u.length))
	if u.next != nil {
		flooredLength += params.SnakeWidth
	}
	if u.direction.isVertical() {
		drawWidth, drawHeight = params.SnakeWidth, flooredLength
	} else {
		drawWidth, drawHeight = flooredLength, params.SnakeWidth
	}

	// Update the options
	u.drawOpts.Uniforms["Direction"] = float32(u.direction)
	u.drawOpts.Uniforms["Size"] = []float32{drawWidth, drawHeight}
	u.drawOpts.Uniforms["ProxToFood"] = proxToFood
}

func (u *unit) moveUp(dist float64) {
	u.headCenter.y -= dist

	// teleport if head center is offscreen.
	if params.TeleportActive && (u.headCenter.y < 0) {
		u.headCenter.y += params.ScreenHeight
	}
}

func (u *unit) moveDown(dist float64) {
	u.headCenter.y += dist

	// teleport if head center is offscreen.
	if params.TeleportActive && (u.headCenter.y > params.ScreenHeight) {
		u.headCenter.y -= params.ScreenHeight
	}
}

func (u *unit) moveRight(dist float64) {
	u.headCenter.x += dist

	// teleport if head center is offscreen.
	if params.TeleportActive && (u.headCenter.x > params.ScreenWidth) {
		u.headCenter.x -= params.ScreenWidth
	}
}

func (u *unit) moveLeft(dist float64) {
	u.headCenter.x -= dist

	// teleport if head center is offscreen.
	if params.TeleportActive && (u.headCenter.x < 0) {
		u.headCenter.x += params.ScreenWidth
	}
}

func (u *unit) markHeadCenters(dst *ebiten.Image) {
	markPoint(dst, u.headCenter, 4, params.ColorFood)

	var offset float64 = 0
	if u.next == nil {
		offset = params.SnakeWidth
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
	markPoint(dst, backCenter, 4, params.ColorFood)
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
	return params.ShaderRound
}

func (u *unit) drawDebugInfo(dst *ebiten.Image) {
	u.markHeadCenters(dst)
	for iRect := range u.rectsDrawable {
		rect := u.rectsDrawable[iRect]
		rect.drawOuterRect(dst, params.ColorFood)
	}
}
