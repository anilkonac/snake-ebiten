package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type directionT uint8
type snakeLengthT uint16

const maxSnakeLength = 500

const (
	directionUp directionT = iota
	directionDown
	directionLeft
	directionRight
	directionTotal
)

type rotation struct {
	centerX       float64
	centerY       float64
	directionFrom directionT
	directionTo   directionT
}

type snake struct {
	speed     uint8
	rotations []rotation
	units     []unit
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

	snake := &snake{
		speed:     speed,
		rotations: make([]rotation, 0, maxSnakeLength),
		units:     make([]unit, 0, maxSnakeLength),
	}

	snake.units = append(snake.units, unit{
		headCenterX: centerX,
		headCenterY: centerY,
		direction:   direction,
		length:      length,
	})

	return snake
}

func (s *snake) update() {

	// Divide snake into units according to rotations
	// numUnits := len(s.rotations) + 1

	// Update units' head position
	travelDistance := float64(s.speed) * deltaTime
	for indexUnit := 0; indexUnit < len(s.units); indexUnit++ {
		curUnit := &s.units[indexUnit]
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

func (s *snake) draw(screen *ebiten.Image) {
	for indexUnit := 0; indexUnit < len(s.units); indexUnit++ {

		curUnit := &s.units[indexUnit]

		// Select color of the unit
		var color color.Color
		if indexUnit%2 == 0 {
			color = colorSnake
		} else {
			color = colorSnake2
		}

		curUnit.draw(screen, color)
	}
}

func (s *snake) headUnit() *unit {
	if len(s.units) == 0 {
		panic("Headless snake!!")
	}

	return &s.units[0]
}

func (s *snake) rotateTo(direction directionT) {
	headUnit := s.headUnit()

	s.rotations = append(s.rotations, rotation{
		centerX:       headUnit.headCenterX,
		centerY:       headUnit.headCenterY,
		directionFrom: headUnit.direction,
		directionTo:   direction,
	})
}
