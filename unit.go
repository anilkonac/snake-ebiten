package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const epsilon = 0.001

type unit struct {
	headCenterX float64
	headCenterY float64
	direction   directionT
	length      float64
	rects       []rectF64 // rectangles that are used for both collision checking and drawing
	color       *color.RGBA
	next        *unit
	prev        *unit
}

func newUnit(headCenterX, headCenterY, length float64, direction directionT, color *color.RGBA) *unit {
	newUnit := &unit{
		headCenterX: headCenterX,
		headCenterY: headCenterY,
		direction:   direction,
		length:      length,
		color:       color,
	}
	newUnit.creteRects()

	return newUnit
}

func (u *unit) moveUp(dist float64) {
	u.headCenterY -= dist

	// teleport if head center is offscreen.
	if u.headCenterY < 0 {
		u.headCenterY += screenHeight
	}
}

func (u *unit) moveDown(dist float64) {
	u.headCenterY += dist

	// teleport if head center is offscreen.
	if u.headCenterY > screenHeight {
		u.headCenterY -= screenHeight
	}
}

func (u *unit) moveRight(dist float64) {
	u.headCenterX += dist

	// teleport if head center is offscreen.
	if u.headCenterX > screenWidth {
		u.headCenterX -= screenWidth
	}
}

func (u *unit) moveLeft(dist float64) {
	u.headCenterX -= dist

	// teleport if head center is offscreen.
	if u.headCenterX < 0 {
		u.headCenterX += screenWidth
	}
}

func (u *unit) draw(screen *ebiten.Image) {
	for _, rect := range u.rects {
		ebitenutil.DrawRect(screen, rect.x, rect.y, rect.width, rect.height, u.color)
		if debugUnits {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%3.3f, %3.3f", rect.x, rect.y), int(rect.x)-90, int(rect.y)-15)
			bottomX := rect.x + rect.width
			bottomY := rect.y + rect.height
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%3.3f, %3.3f", bottomX, bottomY), int(bottomX), int(bottomY))
		}
	}
}

func (u *unit) creteRects() {
	// Create the rectangle to be sliced.
	var pureRect rectF64
	length64 := float64(u.length)
	switch u.direction {
	case directionRight:
		pureRect = rectF64{
			x:      u.headCenterX - length64 + halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  length64,
			height: snakeWidth,
		}
	case directionLeft:
		pureRect = rectF64{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  length64,
			height: snakeWidth,
		}
	case directionUp:
		pureRect = rectF64{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  snakeWidth,
			height: length64,
		}
	case directionDown:
		pureRect = rectF64{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - length64 + halfSnakeWidth,
			width:  snakeWidth,
			height: length64,
		}
	default:
		panic("Wrong unit direction!!")
	}

	u.rects = make([]rectF64, 0, 4) // Remove old rectangles
	pureRect.slice(&u.rects)        // Create rectangles divided at screen edges
}

func (a *unit) intersects(b *unit) bool {
	for _, rectA := range a.rects {
		for _, rectB := range b.rects {
			if intersects(rectA, rectB) {
				return true
			}

			continue
		}
	}

	return false
}
