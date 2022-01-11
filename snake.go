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

type unit struct {
	relX      int
	relY      int
	direction uint8
}

type snake struct {
	headPosX float64
	headPosY float64
	speed    uint8
	units    []unit
}

func newSnake(posX float64, posY float64, direction uint8, speed uint8, length uint8) *snake {
	if direction >= directionTotal {
		panic("direction parameter is invalid.")
	}
	if posX > screenWidth {
		panic("initial x position of the snake is out of screen.")
	}
	if posY > screenHeight {
		panic("initial y position of the snake is out of screen.")
	}

	// Create snake
	snake := &snake{
		headPosX: posX - unitLength/2.0, // set center of the head
		headPosY: posY - unitLength/2.0, // set center of the head
		speed:    speed,
		units:    make([]unit, length, math.MaxUint8),
	}

	// Initialize units of snake
	iLength := int(length)
	for i := 0; i < iLength; i++ {
		curUnit := &snake.units[i]
		curUnit.direction = direction

		// Compute relative position of current unit
		distanceToHead := i * unitLength
		switch direction {
		case directionUp:
			curUnit.relX = 0
			curUnit.relY = distanceToHead
		case directionDown:
			curUnit.relX = 0
			curUnit.relY = -distanceToHead
		case directionRight:
			curUnit.relX = -distanceToHead
			curUnit.relY = 0
		case directionLeft:
			curUnit.relX = distanceToHead
			curUnit.relY = 0
		}
	}

	return snake
}

func (s snake) headDirection() uint8 {
	if len(s.units) == 0 {
		panic("Snake does not have a head!!")
	}

	return s.units[0].direction
}
