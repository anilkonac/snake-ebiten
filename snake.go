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
	halfUnitLength = unitLength / 2.0
)

type unit struct {
	centerX   float64
	centerY   float64
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

func (u *unit) moveUp(dist float64) {
	u.centerY -= dist

	// Teleport if center is off the screen
	if u.centerY < 0 {
		u.centerY += screenWidth
	}
}

func (u *unit) moveDown(dist float64) {
	u.centerY += dist

	// Teleport if center is off the screen
	if u.centerY > screenHeight {
		u.centerY -= screenWidth
	}
}

func (u *unit) moveLeft(dist float64) {
	u.centerX -= dist

	// Teleport if center is off the screen
	if u.centerX < 0 {
		u.centerX += screenWidth
	}
}

func (u *unit) moveRight(dist float64) {
	u.centerX += dist

	// Teleport if center is off the screen
	if u.centerX > screenWidth {
		u.centerX -= screenWidth
	}
}

// Checks if unit should be mirrored
// Returns which axes should be mirrored, and position of center on mirrored axes
func (u unit) checkOffScreen() (x, y bool, locCenterX, locCenterY uint8) {
	if u.centerX-halfUnitLength < 0 {
		x = true
		locCenterX = directionLeft
	} else if u.centerX+halfUnitLength > screenWidth {
		x = true
		locCenterX = directionRight
	}

	if u.centerY-halfUnitLength < 0 {
		y = true
		locCenterY = directionUp
	} else if u.centerY+halfUnitLength > screenHeight {
		y = true
		locCenterY = directionDown
	}

	return
}
