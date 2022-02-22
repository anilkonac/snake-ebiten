package game

import (
	"image/color"
)

const epsilon = 0.001

type unit struct {
	headCenterX float64
	headCenterY float64
	length      float64
	direction   directionT
	rects       []rectF64 // rectangles that are used for both collision checking and drawing
	color       *color.RGBA
	next        *unit
	prev        *unit
}

func newUnit(headCenterX, headCenterY, length float64, direction directionT, color *color.RGBA) *unit {
	newUnit := &unit{
		headCenterX: headCenterX,
		headCenterY: headCenterY,
		length:      length,
		direction:   direction,
		color:       color,
	}
	newUnit.creteRects()

	return newUnit
}

func (u *unit) moveUp(dist float64) {
	u.headCenterY -= dist

	// teleport if head center is offscreen.
	if u.headCenterY < 0 {
		u.headCenterY += ScreenHeight
	}
}

func (u *unit) moveDown(dist float64) {
	u.headCenterY += dist

	// teleport if head center is offscreen.
	if u.headCenterY > ScreenHeight {
		u.headCenterY -= ScreenHeight
	}
}

func (u *unit) moveRight(dist float64) {
	u.headCenterX += dist

	// teleport if head center is offscreen.
	if u.headCenterX > ScreenWidth {
		u.headCenterX -= ScreenWidth
	}
}

func (u *unit) moveLeft(dist float64) {
	u.headCenterX -= dist

	// teleport if head center is offscreen.
	if u.headCenterX < 0 {
		u.headCenterX += ScreenWidth
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
	pureRect.split(&u.rects)        // Create split rectangles on screen edges.
}

// Implement slicer interface
// --------------------------
func (u *unit) slice() []rectF64 {
	return u.rects
}

// Implement collidable interface
// ------------------------------
func (u *unit) collEnabled() bool {
	return true
}

// Implement drawable interface
// ------------------------------
func (u *unit) drawEnabled() bool {
	return true
}

func (u *unit) Color() color.Color {
	return u.color
}
