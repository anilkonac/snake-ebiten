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
	"github.com/anilkonac/snake-ebiten/game/shaders"
	t "github.com/anilkonac/snake-ebiten/game/tools"
	"github.com/hajimehoshi/ebiten/v2"
)

var shaderSnakeHead *ebiten.Shader

func init() {
	shaderSnakeHead = t.NewShader(shaders.SnakeHead)
}

type unit struct {
	headCenter     t.Vec64
	length         float64
	direction      directionT
	rectsCollision []t.RectF32
	rectsDrawable  []t.RectF32
	color          *color.RGBA
	next           *unit
	prev           *unit
	drawOpts       ebiten.DrawTrianglesShaderOptions
}

func newUnit(headCenter t.Vec64, length float64, direction directionT, color *color.RGBA) *unit {
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
	var rectDraw, rectColl *t.RectF32

	rectColl = u.createRectColl()
	rectDraw = u.createRectDraw(rectColl)

	// Remove old rectangles
	u.rectsCollision = make([]t.RectF32, 0, 4)
	if u.next != nil {
		u.rectsDrawable = make([]t.RectF32, 0, 4)
	}

	// Create split rectangles on screen edges.
	rectColl.Split(&u.rectsCollision)
	if u.next == nil {
		u.rectsDrawable = u.rectsCollision
		return
	}
	rectDraw.Split(&u.rectsDrawable)
}

func (u *unit) createRectColl() (rectColl *t.RectF32) {
	length32 := float32(math.Floor(u.length))
	flCenter := u.headCenter.Floor().To32()

	switch u.direction {
	case directionRight:
		rectColl = t.NewRect(
			t.Vec32{
				X: flCenter.X - length32 + params.RadiusSnake,
				Y: flCenter.Y - params.RadiusSnake,
			},
			t.Vec32{X: length32, Y: params.SnakeWidth},
		)
	case directionLeft:
		rectColl = t.NewRect(
			t.Vec32{
				X: flCenter.X - params.RadiusSnake,
				Y: flCenter.Y - params.RadiusSnake,
			},
			t.Vec32{X: length32, Y: params.SnakeWidth},
		)
	case directionUp:
		rectColl = t.NewRect(
			t.Vec32{
				X: flCenter.X - params.RadiusSnake,
				Y: flCenter.Y - params.RadiusSnake,
			},
			t.Vec32{X: params.SnakeWidth, Y: length32},
		)
	case directionDown:
		rectColl = t.NewRect(
			t.Vec32{
				X: flCenter.X - params.RadiusSnake,
				Y: flCenter.Y - length32 + params.RadiusSnake,
			},
			t.Vec32{X: params.SnakeWidth, Y: length32})
	default:
		panic("Wrong unit direction!!")
	}

	return
}

func (u *unit) createRectDraw(rectColl *t.RectF32) (rectDraw *t.RectF32) {
	if u.next == nil {
		rectDraw = rectColl
		return
	}

	switch u.direction {
	case directionRight:
		rectDraw = t.NewRect(t.Vec32{X: rectColl.Pos.X - params.SnakeWidth, Y: rectColl.Pos.Y}, t.Vec32{X: rectColl.Size.X + params.SnakeWidth, Y: rectColl.Size.Y})
	case directionLeft:
		rectDraw = t.NewRect(t.Vec32{X: rectColl.Pos.X, Y: rectColl.Pos.Y}, t.Vec32{X: rectColl.Size.X + params.SnakeWidth, Y: rectColl.Size.Y})
	case directionUp:
		rectDraw = t.NewRect(t.Vec32{X: rectColl.Pos.X, Y: rectColl.Pos.Y}, t.Vec32{X: rectColl.Size.X, Y: rectColl.Size.Y + params.SnakeWidth})
	case directionDown:
		rectDraw = t.NewRect(t.Vec32{X: rectColl.Pos.X, Y: rectColl.Pos.Y - params.SnakeWidth}, t.Vec32{X: rectColl.Size.X, Y: rectColl.Size.Y + params.SnakeWidth})
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
	u.headCenter.Y -= dist

	// teleport if head center is offscreen.
	if params.TeleportActive && (u.headCenter.Y < 0) {
		u.headCenter.Y += params.ScreenHeight
	}
}

func (u *unit) moveDown(dist float64) {
	u.headCenter.Y += dist

	// teleport if head center is offscreen.
	if params.TeleportActive && (u.headCenter.Y > params.ScreenHeight) {
		u.headCenter.Y -= params.ScreenHeight
	}
}

func (u *unit) moveRight(dist float64) {
	u.headCenter.X += dist

	// teleport if head center is offscreen.
	if params.TeleportActive && (u.headCenter.X > params.ScreenWidth) {
		u.headCenter.X -= params.ScreenWidth
	}
}

func (u *unit) moveLeft(dist float64) {
	u.headCenter.X -= dist

	// teleport if head center is offscreen.
	if params.TeleportActive && (u.headCenter.X < 0) {
		u.headCenter.X += params.ScreenWidth
	}
}

func (u *unit) markHeadCenters(dst *ebiten.Image) {
	t.MarkPoint(dst, u.headCenter, 4, params.ColorFood)

	var offset float64 = 0
	if u.next == nil {
		offset = params.SnakeWidth
	}

	backCenter := u.headCenter
	switch u.direction {
	case directionUp:
		backCenter.Y = u.headCenter.Y + u.length - offset
	case directionDown:
		backCenter.Y = u.headCenter.Y - u.length + offset
	case directionRight:
		backCenter.X = u.headCenter.X - u.length + offset
	case directionLeft:
		backCenter.X = u.headCenter.X + u.length - offset
	}
	// mark head center at the other side
	t.MarkPoint(dst, backCenter, 4, params.ColorFood)
}

// Implement collidable interface
// ------------------------------
func (u *unit) collEnabled() bool {
	return true
}

func (u *unit) collisionRects() []t.RectF32 {
	return u.rectsCollision
}

// Implement drawable interface
// ----------------------------
func (u *unit) drawEnabled() bool {
	return true
}

func (u *unit) drawableRects() []t.RectF32 {
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
		rect.DrawOuterRect(dst, params.ColorFood)
	}
}
