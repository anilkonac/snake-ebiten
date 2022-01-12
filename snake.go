package main

import "math"

const (
	directionUp uint8 = iota
	directionDown
	directionLeft
	directionRight
	directionTotal
)

const (
	unitLength   = 20
	centerOffset = unitLength / 2.0
)

type unit struct {
	posX      float64
	posY      float64
	direction uint8
}

type snake struct {
	speed uint8
	units []unit
}

func newSnake(centerX float64, centerY float64, direction uint8, speed uint8, length uint8) *snake {
	if direction >= directionTotal {
		panic("direction parameter is invalid.")
	}
	if centerX > screenWidth {
		panic("initial x position of the snake is out of screen.")
	}
	if centerY > screenHeight {
		panic("initial x position of the snake is out of screen.")
	}

	snake := &snake{speed: speed}

	// Generate units of the snake
	snake.units = make([]unit, length, math.MaxUint8)
	for i := uint8(0); i < length; i++ {
		curUnit := &snake.units[i]
		curUnit.direction = direction

		// Compute position of current unit
		distanceToHead := float64(i) * unitLength
		switch direction {
		case directionUp:
			curUnit.posX = centerX
			curUnit.posY = centerY - distanceToHead
		case directionDown:
			curUnit.posX = centerX
			curUnit.posY = centerY + distanceToHead
		case directionRight:
			curUnit.posX = centerX - distanceToHead
			curUnit.posY = centerY
		case directionLeft:
			curUnit.posX = centerX + distanceToHead
			curUnit.posY = centerY
		}
	}

	return snake
}

func (s *snake) update() {
	// Update snake position
	for indexUnit := 0; indexUnit < len(s.units); indexUnit++ {
		curUnit := &s.units[indexUnit]

		travelDistance := float64(s.speed) * deltaTime
		switch curUnit.direction {
		case directionRight:
			curUnit.moveRight(travelDistance)
		case directionLeft:
			curUnit.moveLeft(travelDistance)
		case directionUp:
			curUnit.moveUp(travelDistance)
		case directionDown:
			curUnit.moveDown(travelDistance)
		}
	}
}

func (u *unit) moveUp(dist float64) {
	u.posY -= dist
	if u.posY < 0 {
		u.posY += screenWidth
	}
}

func (u *unit) moveDown(dist float64) {
	u.posY += dist
	if u.posY > screenHeight {
		u.posY -= screenWidth
	}
}

func (u *unit) moveLeft(dist float64) {
	u.posX -= dist
	if u.posX < 0 {
		u.posX += screenWidth
	}
}

func (u *unit) moveRight(dist float64) {
	u.posX += dist
	if u.posX > screenWidth {
		u.posX -= screenWidth
	}
}
