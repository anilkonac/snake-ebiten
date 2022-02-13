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
	rects       []rectF64 // rectangles that are used for both collision checking and drawing
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
	for _, rect := range u.rects {
		rect.draw(screen, u.color)
	}
}

func (u *unit) creteRects() {
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

	u.rects = make([]rectF64, 0, 4) // Reset rectangles of this unit
	pureRect.slice(&u.rects)        // Create rectangles divided at screen edges

}

func (a *unit) intersects(b *unit) bool {
	for _, rectA := range a.rects {
		for _, rectB := range b.rects {
			aRightX := rectA.x + rectA.width
			bRightX := rectB.x + rectB.width
			aBottomY := rectA.y + rectA.height
			bBottomY := rectB.y + rectB.height

			if (rectA.x-rectB.x <= epsilon) && (aRightX-rectB.x <= epsilon) { // a is on the left side of b
				continue
			}

			if (rectA.x-bRightX >= -epsilon) && (aRightX-bRightX >= -epsilon) { // a is on the right side of b
				continue
			}

			if (rectA.y-rectB.y <= epsilon) && (aBottomY-rectB.y <= epsilon) { // a is above b
				continue
			}

			if (rectA.y-bBottomY >= -epsilon) && (aBottomY-bBottomY >= -epsilon) { // a is under b
				continue
			}

			return true
		}
	}

	return false
}
