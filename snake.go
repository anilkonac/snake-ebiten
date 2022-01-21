package main

import "github.com/hajimehoshi/ebiten/v2"

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
	posX      float64
	posY      float64
	direction directionT
}

type snake struct {
	headCenterX float64
	headCenterY float64
	speed       uint8
	direction   directionT
	length      snakeLengthT
	rotations   []rotation
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
		headCenterX: centerX,
		headCenterY: centerY,
		speed:       speed,
		direction:   direction,
		length:      length,
		rotations:   make([]rotation, 0, maxSnakeLength),
	}

	return snake
}

func (s *snake) update() {
	// Update snake's head position
	travelDistance := float64(s.speed) * deltaTime
	switch s.direction {
	case directionRight:
		s.moveRight(travelDistance)
	case directionLeft:
		s.moveLeft(travelDistance)
	case directionUp:
		s.moveUp(travelDistance)
	case directionDown:
		s.moveDown(travelDistance)
	}
}

func (s *snake) draw(screen *ebiten.Image) {
	snakeLength64 := float64(s.length)
	switch s.direction {
	// Create a screenRect whose x and y coordinates are top left corner. Then draw it.
	case directionRight:
		screenRect{
			x:      s.headCenterX - unitLength*snakeLength64 + halfUnitLength,
			y:      s.headCenterY - halfUnitLength,
			width:  unitLength * snakeLength64,
			height: unitLength,
		}.draw(screen, colorSnake)
	case directionLeft:
		screenRect{
			x:      s.headCenterX - halfUnitLength,
			y:      s.headCenterY - halfUnitLength,
			width:  unitLength * snakeLength64,
			height: unitLength,
		}.draw(screen, colorSnake)
	case directionUp:
		screenRect{
			x:      s.headCenterX - halfUnitLength,
			y:      s.headCenterY - halfUnitLength,
			width:  unitLength,
			height: unitLength * snakeLength64,
		}.draw(screen, colorSnake)
	case directionDown:
		screenRect{
			x:      s.headCenterX - halfUnitLength,
			y:      s.headCenterY - unitLength*snakeLength64 + halfUnitLength,
			width:  unitLength,
			height: unitLength * snakeLength64,
		}.draw(screen, colorSnake)
	}
}

func (s *snake) moveUp(dist float64) {
	s.headCenterY -= dist

	// teleport if head center is offscreen.
	if s.headCenterY < 0 {
		s.headCenterY += screenHeight
	}
}

func (s *snake) moveDown(dist float64) {
	s.headCenterY += dist

	// teleport if head center is offscreen.
	if s.headCenterY > screenHeight {
		s.headCenterY -= screenHeight
	}
}

func (s *snake) moveRight(dist float64) {
	s.headCenterX += dist

	// teleport if head center is offscreen.
	if s.headCenterX > screenWidth {
		s.headCenterX -= screenWidth
	}
}

func (s *snake) moveLeft(dist float64) {
	s.headCenterX -= dist

	// teleport if head center is offscreen.
	if s.headCenterX < 0 {
		s.headCenterX += screenWidth
	}
}
