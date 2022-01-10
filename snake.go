package main

import "math"

const (
	directionUp uint8 = iota
	directionDown
	directionLeft
	directionRight
	directionTotal
)

const unitLength = 20

type Unit struct {
	posX      float64
	posY      float64
	direction uint8
}

type Snake struct {
	speed uint8
	units []Unit
}

func newSnake(posX float64, posY float64, direction uint8, speed uint8, length uint8) *Snake {
	if direction >= directionTotal {
		panic("direction parameter is invalid.")
	}
	if posX > screenWidth {
		panic("initial x position of the snake is out of screen.")
	}
	if posY > screenHeight {
		panic("initial x position of the snake is out of screen.")
	}

	snake := &Snake{speed: speed}

	// Generate units of the snake
	snake.units = make([]Unit, length, math.MaxUint8)
	iLength := int(length)
	for i := 0; i < iLength; i++ {
		curUnit := &snake.units[i]
		curUnit.direction = direction

		// Compute position of current unit
		distanceToHead := float64(i * unitLength)
		switch direction {
		case directionUp:
			curUnit.posX = posX
			curUnit.posY = posY - distanceToHead
		case directionDown:
			curUnit.posX = posX
			curUnit.posY = posY + distanceToHead
		case directionRight:
			curUnit.posX = posX - distanceToHead
			curUnit.posY = posY
		case directionLeft:
			curUnit.posX = posX + distanceToHead
			curUnit.posY = posY
		}
	}

	return snake
}
