package main

import "math"

const (
	directionUp uint8 = iota
	directionDown
	directionLeft
	directionRight
	directionTotal
)

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
			curUnit.centerX = centerX
			curUnit.centerY = centerY - distanceToHead
		case directionDown:
			curUnit.centerX = centerX
			curUnit.centerY = centerY + distanceToHead
		case directionRight:
			curUnit.centerX = centerX - distanceToHead
			curUnit.centerY = centerY
		case directionLeft:
			curUnit.centerX = centerX + distanceToHead
			curUnit.centerY = centerY
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
