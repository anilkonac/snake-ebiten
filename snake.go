package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type directionT uint8
type snakeLengthT uint16

const (
	directionUp directionT = iota
	directionDown
	directionLeft
	directionRight
	directionTotal
)

type snake struct {
	speed    uint8
	unitHead *unit
	unitTail *unit
}

func newSnake(centerX float64, centerY float64, direction directionT, speed uint8, length snakeLengthT) *snake {
	if direction >= directionTotal {
		panic("direction parameter is invalid.")
	}
	if centerX > screenWidth {
		panic("Initial x position of the snake is out of the screen.")
	}
	if centerY > screenHeight {
		panic("Initial y position of the snake is out of the screen.")
	}

	initialUnit := &unit{
		headCenterX: centerX,
		headCenterY: centerY,
		direction:   direction,
		length_px:   float64(length) * snakeWidth,
		next:        nil,
		prev:        nil,
	}

	snake := &snake{
		speed:    speed,
		unitHead: initialUnit,
		unitTail: initialUnit,
	}

	return snake
}

func (s *snake) update() {
	travelDistance := float64(s.speed) * deltaTime

	// Increse head length
	s.unitHead.length_px += travelDistance

	// Decrease tail length
	s.unitTail.length_px -= travelDistance

	// Rotate tail if its length is less than width of the snake
	if s.unitTail.length_px <= snakeWidth {
		s.unitTail.direction = s.unitTail.prev.direction
	}

	// Destroy tail unit if its length is not positive
	if s.unitTail.length_px <= 0 {
		s.unitTail = s.unitTail.prev
		s.unitTail.next = nil
	}

	// Move head
	switch s.unitHead.direction {
	case directionRight:
		s.unitHead.moveRight(travelDistance)
	case directionLeft:
		s.unitHead.moveLeft(travelDistance)
	case directionUp:
		s.unitHead.moveUp(travelDistance)
	case directionDown:
		s.unitHead.moveDown(travelDistance)
	}
}

func (s *snake) draw(screen *ebiten.Image) {

	curUnit := s.unitHead
	indexUnit := 0
	for curUnit != nil {
		// Select color of the unit
		var color color.Color
		if indexUnit%2 == 0 {
			color = colorSnake
		} else {
			color = colorSnake2
		}

		curUnit.draw(screen, color)
		curUnit = curUnit.next
		indexUnit++
	}
}

func (s *snake) rotateTo(direction directionT) {

	// Add a new unit whose direction is the passed parameter to the head
	oldHead := s.unitHead
	newHead := &unit{
		headCenterX: oldHead.headCenterX,
		headCenterY: oldHead.headCenterY,
		direction:   direction,
		length_px:   0,
		next:        oldHead,
		prev:        nil,
	}
	oldHead.prev = newHead
	s.unitHead = newHead
}
