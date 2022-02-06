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
	rect        *screenRect
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
	if u.rect == nil {
		u.creteRect()
	}

	u.rect.draw(screen, u.color)
}

func (u *unit) creteRect() {
	length64 := float64(u.length)
	switch u.direction {
	case directionRight:
		u.rect = &screenRect{
			x:      u.headCenterX - length64 + halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  length64,
			height: snakeWidth,
		}
	case directionLeft:
		u.rect = &screenRect{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  length64,
			height: snakeWidth,
		}
	case directionUp:
		u.rect = &screenRect{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  snakeWidth,
			height: length64,
		}
	case directionDown:
		u.rect = &screenRect{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - length64 + halfSnakeWidth,
			width:  snakeWidth,
			height: length64,
		}
	}
}
