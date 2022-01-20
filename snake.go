package main

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
