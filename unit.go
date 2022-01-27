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

func (u *unit) draw(screen *ebiten.Image, color color.Color) {
	length64 := float64(u.length)
	switch u.direction {
	// Create a screenRect whose x and y coordinates are top left corner. Then draw it.
	case directionRight:
		screenRect{
			x:      u.headCenterX - length64 + halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  length64,
			height: snakeWidth,
		}.draw(screen, color)
	case directionLeft:
		screenRect{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  length64,
			height: snakeWidth,
		}.draw(screen, color)
	case directionUp:
		screenRect{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - halfSnakeWidth,
			width:  snakeWidth,
			height: length64,
		}.draw(screen, color)
	case directionDown:
		screenRect{
			x:      u.headCenterX - halfSnakeWidth,
			y:      u.headCenterY - length64 + halfSnakeWidth,
			width:  snakeWidth,
			height: length64,
		}.draw(screen, color)
	}
}
