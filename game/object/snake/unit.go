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

package snake

import (
	"image/color"
	"math"

	c "github.com/anilkonac/snake-ebiten/game/core"
	"github.com/anilkonac/snake-ebiten/game/param"
	"github.com/anilkonac/snake-ebiten/game/shader"
	"github.com/hajimehoshi/ebiten/v2"
)

var shaderSnakeHead *ebiten.Shader

func init() {
	shaderSnakeHead = shader.New(shader.PathSnakeHead)
}

type Unit struct {
	HeadCenter      c.Vec64
	length          float64
	Direction       DirectionT
	CompCollision   c.TeleComp
	CompBody        c.TeleComp
	CompTriangDebug c.TeleCompTriang
	CompTriangHead  c.TeleCompTriang
	CompTriangTail  c.TeleCompTriang
	Next            *Unit
	prev            *Unit
	drawOpts        ebiten.DrawTrianglesShaderOptions
}

func NewUnit(headCenter c.Vec64, length float64, direction DirectionT, color *color.RGBA) *Unit {
	newUnit := &Unit{
		HeadCenter: headCenter,
		length:     length,
		Direction:  direction,
		drawOpts: ebiten.DrawTrianglesShaderOptions{
			Uniforms: map[string]interface{}{
				"RadiusMouth": float32(param.RadiusMouth),
			},
		},
	}
	newUnit.SetColor(color)
	newUnit.update(param.MouthAnimStartDistance)

	return newUnit
}

func (u *Unit) updateRects() {
	// Create rectangles for drawing and collision. They are going to split.
	var rectDraw, rectColl *c.RectF32

	rectColl = u.createRectCollision()
	rectDraw = u.createRectDraw(rectColl)
	rectDrawHead := u.createRectHead()
	rectDrawBody := u.createRectBody(rectColl)

	u.CompCollision.Update(rectColl)
	u.CompTriangDebug.Update(rectDraw)
	u.CompTriangHead.Update(rectDrawHead)
	u.CompBody.Update(rectDrawBody)

	// If current unit is the tail unit
	if u.Next == nil {
		rectDrawTail := u.createRectTail(rectDrawHead)
		u.CompTriangTail.Update(rectDrawTail)
	}
}

func (u *Unit) createRectCollision() (rectColl *c.RectF32) {
	length32 := float32(math.Floor(u.length))
	flCenter := u.HeadCenter.Floor().To32()

	switch u.Direction {
	case DirectionRight:
		rectColl = c.NewRect(
			c.Vec32{
				X: flCenter.X - length32 + param.RadiusSnake,
				Y: flCenter.Y - param.RadiusSnake,
			},
			c.Vec32{X: length32, Y: param.SnakeWidth},
		)
	case DirectionLeft:
		rectColl = c.NewRect(
			c.Vec32{
				X: flCenter.X - param.RadiusSnake,
				Y: flCenter.Y - param.RadiusSnake,
			},
			c.Vec32{X: length32, Y: param.SnakeWidth},
		)
	case DirectionUp:
		rectColl = c.NewRect(
			c.Vec32{
				X: flCenter.X - param.RadiusSnake,
				Y: flCenter.Y - param.RadiusSnake,
			},
			c.Vec32{X: param.SnakeWidth, Y: length32},
		)
	case DirectionDown:
		rectColl = c.NewRect(
			c.Vec32{
				X: flCenter.X - param.RadiusSnake,
				Y: flCenter.Y - length32 + param.RadiusSnake,
			},
			c.Vec32{X: param.SnakeWidth, Y: length32})
	default:
		panic("Wrong unit direction!!")
	}

	return
}

func (u *Unit) createRectDraw(rectColl *c.RectF32) (rectDraw *c.RectF32) {
	if u.Next == nil {
		rectDraw = rectColl
		return
	}

	switch u.Direction {
	case DirectionRight:
		rectDraw = c.NewRect(c.Vec32{X: rectColl.Pos.X - param.SnakeWidth, Y: rectColl.Pos.Y}, c.Vec32{X: rectColl.Size.X + param.SnakeWidth, Y: rectColl.Size.Y})
	case DirectionLeft:
		rectDraw = c.NewRect(c.Vec32{X: rectColl.Pos.X, Y: rectColl.Pos.Y}, c.Vec32{X: rectColl.Size.X + param.SnakeWidth, Y: rectColl.Size.Y})
	case DirectionUp:
		rectDraw = c.NewRect(c.Vec32{X: rectColl.Pos.X, Y: rectColl.Pos.Y}, c.Vec32{X: rectColl.Size.X, Y: rectColl.Size.Y + param.SnakeWidth})
	case DirectionDown:
		rectDraw = c.NewRect(c.Vec32{X: rectColl.Pos.X, Y: rectColl.Pos.Y - param.SnakeWidth}, c.Vec32{X: rectColl.Size.X, Y: rectColl.Size.Y + param.SnakeWidth})
	default:
		panic("Wrong unit direction!!")
	}

	return
}

func (u *Unit) createRectHead() *c.RectF32 {
	headCenter32 := u.HeadCenter.Floor().To32()
	return c.NewRect(c.Vec32{X: headCenter32.X - param.RadiusSnake, Y: headCenter32.Y - param.RadiusSnake}, c.Vec32{X: param.SnakeWidth, Y: param.SnakeWidth})
}

func (u *Unit) createRectTail(rectHead *c.RectF32) (rectTail *c.RectF32) {
	size := c.Vec32{X: param.SnakeWidth, Y: param.SnakeWidth}
	switch u.Direction {
	case DirectionUp:
		rectTail = c.NewRect(c.Vec32{X: rectHead.Pos.X, Y: rectHead.Pos.Y + float32(u.length) - param.SnakeWidth}, size)
	case DirectionDown:
		rectTail = c.NewRect(c.Vec32{X: rectHead.Pos.X, Y: rectHead.Pos.Y - float32(u.length) + param.SnakeWidth}, size)
	case DirectionLeft:
		rectTail = c.NewRect(c.Vec32{X: rectHead.Pos.X + float32(u.length) - param.SnakeWidth, Y: rectHead.Pos.Y}, size)
	case DirectionRight:
		rectTail = c.NewRect(c.Vec32{X: rectHead.Pos.X - float32(u.length) + param.SnakeWidth, Y: rectHead.Pos.Y}, size)
	}
	return
}

func (u *Unit) createRectBody(rectColl *c.RectF32) (rectBody *c.RectF32) {

	switch u.Direction {
	case DirectionUp:
		rectBody = c.NewRect(
			c.Vec32{X: rectColl.Pos.X, Y: rectColl.Pos.Y + param.RadiusSnake},
			c.Vec32{X: rectColl.Size.X, Y: rectColl.Size.Y - param.SnakeWidth},
		)
		if u.Next != nil {
			rectBody.Size.Y += param.SnakeWidth
		}
	case DirectionDown:
		rectBody = c.NewRect(
			c.Vec32{X: rectColl.Pos.X, Y: rectColl.Pos.Y + param.RadiusSnake},
			c.Vec32{X: rectColl.Size.X, Y: rectColl.Size.Y - param.SnakeWidth},
		)
		if u.Next != nil {
			rectBody.Pos.Y -= param.SnakeWidth
			rectBody.Size.Y += param.SnakeWidth
		}
	case DirectionLeft:
		rectBody = c.NewRect(
			c.Vec32{X: rectColl.Pos.X + param.RadiusSnake, Y: rectColl.Pos.Y},
			c.Vec32{X: rectColl.Size.X - param.SnakeWidth, Y: rectColl.Size.Y},
		)
		if u.Next != nil {
			rectBody.Size.X += param.SnakeWidth
		}
	case DirectionRight:
		rectBody = c.NewRect(
			c.Vec32{X: rectColl.Pos.X + param.RadiusSnake, Y: rectColl.Pos.Y},
			c.Vec32{X: rectColl.Size.X - param.SnakeWidth, Y: rectColl.Size.Y},
		)
		if u.Next != nil {
			rectBody.Size.X += param.SnakeWidth
			rectBody.Pos.X -= param.SnakeWidth
		}
	}
	return
}

func (u *Unit) update(distToFood float32) {
	u.updateRects()
	u.updateDrawOptions(distToFood)
}

func (u *Unit) updateDrawOptions(distToFood float32) {
	// Distance to food
	proxToFood := 1.0 - distToFood/param.MouthAnimStartDistance

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
	if param.TeleportEnabled && (u.HeadCenter.Y < 0) {
		u.HeadCenter.Y += param.ScreenHeight
	}
}

func (u *Unit) moveDown(dist float64) {
	u.HeadCenter.Y += dist

	// teleport if head center is offscreen.
	if param.TeleportEnabled && (u.HeadCenter.Y > param.ScreenHeight) {
		u.HeadCenter.Y -= param.ScreenHeight
	}
}

func (u *Unit) moveRight(dist float64) {
	u.HeadCenter.X += dist

	// teleport if head center is offscreen.
	if param.TeleportEnabled && (u.HeadCenter.X > param.ScreenWidth) {
		u.HeadCenter.X -= param.ScreenWidth
	}
}

func (u *Unit) moveLeft(dist float64) {
	u.HeadCenter.X -= dist

	// teleport if head center is offscreen.
	if param.TeleportEnabled && (u.HeadCenter.X < 0) {
		u.HeadCenter.X += param.ScreenWidth
	}
}

func (u *Unit) markHeadCenters(dst *ebiten.Image) {
	c.MarkPoint(dst, u.HeadCenter, 4, param.ColorFood)

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
	c.MarkPoint(dst, backCenter, 4, param.ColorFood)
}

func (u *Unit) SetColor(clr *color.RGBA) {
	u.CompTriangDebug.SetColor(clr)
	u.CompTriangHead.SetColor(clr)
	u.CompTriangTail.SetColor(clr)
}

// Implement collidable interface
// ------------------------------
func (u *Unit) CollEnabled() bool {
	return true
}

func (u *Unit) CollisionRects() []c.RectF32 {
	return u.CompCollision.Rects[:]
}

// Implement drawable interface
// ----------------------------
func (u *Unit) DrawEnabled() bool {
	return true
}

func (u *Unit) Triangles() ([]ebiten.Vertex, []uint16) {
	return u.CompTriangDebug.Triangles()
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
	for iRect := uint8(0); iRect < u.CompTriangDebug.NumRects; iRect++ {
		u.CompTriangDebug.Rects[iRect].DrawOuterRect(dst, param.ColorFood)
	}
}
