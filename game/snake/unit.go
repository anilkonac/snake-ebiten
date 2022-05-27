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

package snake

import (
	"image/color"
	"math"

	"github.com/anilkonac/snake-ebiten/game/param"
	"github.com/anilkonac/snake-ebiten/game/shader"
	t "github.com/anilkonac/snake-ebiten/game/tool"
	"github.com/hajimehoshi/ebiten/v2"
)

var shaderSnakeHead *ebiten.Shader

func init() {
	shaderSnakeHead = t.NewShader(shader.SnakeHead)
}

type Unit struct {
	HeadCenter     t.Vec64
	length         float64
	Direction      DirectionT
	RectsCollision []t.RectF32
	rectsDrawable  []t.RectF32
	color          *color.RGBA
	Next           *Unit
	prev           *Unit
	drawOpts       ebiten.DrawTrianglesShaderOptions
}

func NewUnit(headCenter t.Vec64, length float64, direction DirectionT, color *color.RGBA) *Unit {
	newUnit := &Unit{
		HeadCenter: headCenter,
		length:     length,
		Direction:  direction,
		color:      color,
		drawOpts: ebiten.DrawTrianglesShaderOptions{
			Uniforms: map[string]interface{}{
				"Radius":      float32(param.RadiusSnake),
				"RadiusMouth": float32(param.RadiusMouth),
			},
		},
	}
	newUnit.update(param.EatingAnimStartDistance)

	return newUnit
}

func (u *Unit) CreateRects() {
	// Create rectangles for drawing and collision. They are going to split.
	var rectDraw, rectColl *t.RectF32

	rectColl = u.createRectColl()
	rectDraw = u.createRectDraw(rectColl)

	// Remove old rectangles
	u.RectsCollision = make([]t.RectF32, 0, 4)
	if u.Next != nil {
		u.rectsDrawable = make([]t.RectF32, 0, 4)
	}

	// Create split rectangles on screen edges.
	rectColl.Split(&u.RectsCollision)
	if u.Next == nil {
		u.rectsDrawable = u.RectsCollision
		return
	}
	rectDraw.Split(&u.rectsDrawable)
}

func (u *Unit) createRectColl() (rectColl *t.RectF32) {
	length32 := float32(math.Floor(u.length))
	flCenter := u.HeadCenter.Floor().To32()

	switch u.Direction {
	case DirectionRight:
		rectColl = t.NewRect(
			t.Vec32{
				X: flCenter.X - length32 + param.RadiusSnake,
				Y: flCenter.Y - param.RadiusSnake,
			},
			t.Vec32{X: length32, Y: param.SnakeWidth},
		)
	case DirectionLeft:
		rectColl = t.NewRect(
			t.Vec32{
				X: flCenter.X - param.RadiusSnake,
				Y: flCenter.Y - param.RadiusSnake,
			},
			t.Vec32{X: length32, Y: param.SnakeWidth},
		)
	case DirectionUp:
		rectColl = t.NewRect(
			t.Vec32{
				X: flCenter.X - param.RadiusSnake,
				Y: flCenter.Y - param.RadiusSnake,
			},
			t.Vec32{X: param.SnakeWidth, Y: length32},
		)
	case DirectionDown:
		rectColl = t.NewRect(
			t.Vec32{
				X: flCenter.X - param.RadiusSnake,
				Y: flCenter.Y - length32 + param.RadiusSnake,
			},
			t.Vec32{X: param.SnakeWidth, Y: length32})
	default:
		panic("Wrong unit direction!!")
	}

	return
}

func (u *Unit) createRectDraw(rectColl *t.RectF32) (rectDraw *t.RectF32) {
	if u.Next == nil {
		rectDraw = rectColl
		return
	}

	switch u.Direction {
	case DirectionRight:
		rectDraw = t.NewRect(t.Vec32{X: rectColl.Pos.X - param.SnakeWidth, Y: rectColl.Pos.Y}, t.Vec32{X: rectColl.Size.X + param.SnakeWidth, Y: rectColl.Size.Y})
	case DirectionLeft:
		rectDraw = t.NewRect(t.Vec32{X: rectColl.Pos.X, Y: rectColl.Pos.Y}, t.Vec32{X: rectColl.Size.X + param.SnakeWidth, Y: rectColl.Size.Y})
	case DirectionUp:
		rectDraw = t.NewRect(t.Vec32{X: rectColl.Pos.X, Y: rectColl.Pos.Y}, t.Vec32{X: rectColl.Size.X, Y: rectColl.Size.Y + param.SnakeWidth})
	case DirectionDown:
		rectDraw = t.NewRect(t.Vec32{X: rectColl.Pos.X, Y: rectColl.Pos.Y - param.SnakeWidth}, t.Vec32{X: rectColl.Size.X, Y: rectColl.Size.Y + param.SnakeWidth})
	default:
		panic("Wrong unit direction!!")
	}

	return
}

func (u *Unit) update(distToFood float32) {
	u.CreateRects() // Update rectangles of this unit
	u.updateDrawOptions(distToFood)
}

func (u *Unit) updateDrawOptions(distToFood float32) {
	// Distance to food
	proxToFood := 1.0 - distToFood/param.EatingAnimStartDistance

	// Specify Size uniform variable
	var drawWidth, drawHeight float32
	flooredLength := float32(math.Floor(u.length))
	if u.Next != nil {
		flooredLength += param.SnakeWidth
	}
	if u.Direction.IsVertical() {
		drawWidth, drawHeight = param.SnakeWidth, flooredLength
	} else {
		drawWidth, drawHeight = flooredLength, param.SnakeWidth
	}

	// Update the options
	u.drawOpts.Uniforms["Direction"] = float32(u.Direction)
	u.drawOpts.Uniforms["Size"] = []float32{drawWidth, drawHeight}
	u.drawOpts.Uniforms["ProxToFood"] = proxToFood
}

func (u *Unit) moveUp(dist float64) {
	u.HeadCenter.Y -= dist

	// teleport if head center is offscreen.
	if param.TeleportActive && (u.HeadCenter.Y < 0) {
		u.HeadCenter.Y += param.ScreenHeight
	}
}

func (u *Unit) moveDown(dist float64) {
	u.HeadCenter.Y += dist

	// teleport if head center is offscreen.
	if param.TeleportActive && (u.HeadCenter.Y > param.ScreenHeight) {
		u.HeadCenter.Y -= param.ScreenHeight
	}
}

func (u *Unit) moveRight(dist float64) {
	u.HeadCenter.X += dist

	// teleport if head center is offscreen.
	if param.TeleportActive && (u.HeadCenter.X > param.ScreenWidth) {
		u.HeadCenter.X -= param.ScreenWidth
	}
}

func (u *Unit) moveLeft(dist float64) {
	u.HeadCenter.X -= dist

	// teleport if head center is offscreen.
	if param.TeleportActive && (u.HeadCenter.X < 0) {
		u.HeadCenter.X += param.ScreenWidth
	}
}

func (u *Unit) markHeadCenters(dst *ebiten.Image) {
	t.MarkPoint(dst, u.HeadCenter, 4, param.ColorFood)

	var offset float64 = 0
	if u.Next == nil {
		offset = param.SnakeWidth
	}

	backCenter := u.HeadCenter
	switch u.Direction {
	case DirectionUp:
		backCenter.Y = u.HeadCenter.Y + u.length - offset
	case DirectionDown:
		backCenter.Y = u.HeadCenter.Y - u.length + offset
	case DirectionRight:
		backCenter.X = u.HeadCenter.X - u.length + offset
	case DirectionLeft:
		backCenter.X = u.HeadCenter.X + u.length - offset
	}
	// mark head center at the other side
	t.MarkPoint(dst, backCenter, 4, param.ColorFood)
}

func (u *Unit) SetColor(clr *color.RGBA) {
	u.color = clr
}

// Implement collidable interface
// ------------------------------
func (u *Unit) CollEnabled() bool {
	return true
}

func (u *Unit) CollisionRects() []t.RectF32 {
	return u.RectsCollision
}

// Implement drawable interface
// ----------------------------
func (u *Unit) DrawEnabled() bool {
	return true
}

func (u *Unit) DrawableRects() []t.RectF32 {
	return u.rectsDrawable
}

func (u *Unit) Color() *color.RGBA {
	return u.color
}

func (u *Unit) DrawOptions() *ebiten.DrawTrianglesShaderOptions {
	return &u.drawOpts
}

func (u *Unit) Shader() *ebiten.Shader {
	if u.prev == nil {
		return shaderSnakeHead
	}
	return param.ShaderRound
}

func (u *Unit) DrawDebugInfo(dst *ebiten.Image) {
	u.markHeadCenters(dst)
	for iRect := range u.rectsDrawable {
		rect := u.rectsDrawable[iRect]
		rect.DrawOuterRect(dst, param.ColorFood)
	}
}
