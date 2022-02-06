package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type unit struct {
	headCenterX float64
	headCenterY float64
	direction   directionT
	length      float64
	color       *color.RGBA
	next        *unit
	prev        *unit
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

// Compute this unit's rectangle parameters according to the direction and draw a rectangle with them.
func (u *unit) draw(screen *ebiten.Image) {
	length64 := float64(u.length)
	switch u.direction {
	case directionRight:
		drawRect(
			screen,
			u.headCenterX-length64+halfSnakeWidth, // x
			u.headCenterY-halfSnakeWidth,          // y
			length64, snakeWidth,                  // width, height
			u.color,
		)
	case directionLeft:
		drawRect(
			screen,
			u.headCenterX-halfSnakeWidth,
			u.headCenterY-halfSnakeWidth,
			length64, snakeWidth,
			u.color,
		)
	case directionUp:
		drawRect(
			screen,
			u.headCenterX-halfSnakeWidth,
			u.headCenterY-halfSnakeWidth,
			snakeWidth, length64,
			u.color,
		)
	case directionDown:
		drawRect(
			screen,
			u.headCenterX-halfSnakeWidth,
			u.headCenterY-length64+halfSnakeWidth,
			snakeWidth, length64,
			u.color,
		)
	}
}
