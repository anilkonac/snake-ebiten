package main

import (
	"math/rand"
	"time"

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
	speed         uint8
	unitHead      *unit
	unitTail      *unit
	turnPrev      *turn
	turnQueue     []*turn
	distAfterTurn float64
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newSnake(centerX, centerY float64, direction directionT, speed uint8, snakeLength snakeLengthT) *snake {
	if direction >= directionTotal {
		panic("direction parameter is invalid.")
	}
	if centerX > screenWidth {
		panic("Initial x position of the snake is off-screen.")
	}
	if centerY > screenHeight {
		panic("Initial y position of the snake is off-screen.")
	}

	initialUnit := &unit{
		headCenterX: centerX,
		headCenterY: centerY,
		direction:   direction,
		length:      float64(snakeLength),
		color:       &colorSnake1,
	}

	initialUnit.creteRect()

	snake := &snake{
		speed:    speed,
		unitHead: initialUnit,
		unitTail: initialUnit,
	}

	return snake
}

func newSnakeRandDir(centerX, centerY float64, speed uint8, snakeLength snakeLengthT) *snake {
	var direction directionT = directionT(rand.Intn(int(directionTotal)))
	return newSnake(centerX, centerY, direction, speed, snakeLength)
}

func (s *snake) update() {
	moveDistance := float64(s.speed) * deltaTime

	// Take the next turn in the queue.
	if len(s.turnQueue) > 0 && s.distAfterTurn > snakeWidth {
		var nextTurn *turn
		nextTurn, s.turnQueue = s.turnQueue[0], s.turnQueue[1:] // Pop front
		s.turnTo(nextTurn, true)
	}

	s.updateHead(moveDistance)
	s.updateTail(moveDistance)
}

func (s *snake) updateHead(dist float64) {
	// Increse head length
	s.unitHead.length += dist

	// Move head
	switch s.unitHead.direction {
	case directionRight:
		s.unitHead.moveRight(dist)
	case directionLeft:
		s.unitHead.moveLeft(dist)
	case directionUp:
		s.unitHead.moveUp(dist)
	case directionDown:
		s.unitHead.moveDown(dist)
	}

	if s.unitHead != s.unitTail {
		s.unitHead.creteRect() // Lazy rect updating
	}

	s.distAfterTurn += dist
}

func (s *snake) updateTail(dist float64) {
	// Decrease tail length
	s.unitTail.length -= dist

	// Rotate tail if its length is less than width of the snake
	if s.unitTail.length <= snakeWidth && s.unitTail.prev != nil {
		s.unitTail.direction = s.unitTail.prev.direction
	}

	// Destroy tail unit if its length is not positive
	if s.unitTail.length <= 0 {
		s.unitTail = s.unitTail.prev
		s.unitTail.next = nil
	}

	s.unitTail.creteRect() // Lazy rect updating

}

func (s *snake) draw(screen *ebiten.Image) {
	curUnit := s.unitHead
	for curUnit != nil {
		curUnit.draw(screen)
		curUnit = curUnit.next
	}
}

func (s *snake) turnTo(newTurn *turn, isFromQueue bool) {
	if !isFromQueue {
		// Check if the new turn is dangerous (twice same turns rapidly).
		if (s.turnPrev != nil) &&
			(s.turnPrev.isTurningLeft == newTurn.isTurningLeft) &&
			(s.distAfterTurn <= snakeWidth) {
			// New turn cannot be taken now, push it into the queue
			s.turnQueue = append(s.turnQueue, newTurn)
			return
		}
		// If there are turns in the queue then add the new turn also into it.
		if len(s.turnQueue) > 0 {
			s.turnQueue = append(s.turnQueue, newTurn)
			return
		}
	}
	s.distAfterTurn = 0

	oldHead := s.unitHead

	// Decide color of new head unit
	newColor := &colorSnake1
	if debugUnits && (oldHead.color == &colorSnake1) {
		newColor = &colorSnake2
	}

	// Create new head unit
	newHead := &unit{
		headCenterX: oldHead.headCenterX,
		headCenterY: oldHead.headCenterY,
		direction:   newTurn.directionTo,
		length:      0,
		color:       newColor,
		next:        oldHead,
	}
	newHead.creteRect()

	// Add the new head unit to the beginning of the unit doubly linked list.
	oldHead.prev = newHead
	s.unitHead = newHead

	s.turnPrev = newTurn
}

func (s *snake) checkIntersection() bool {
	if s.unitHead.next == nil {
		return false
	}

	curUnit := s.unitHead.next.next
	for curUnit != nil {
		if s.unitHead.rect.intersects(curUnit.rect) {
			return true
		}
		curUnit = curUnit.next
	}
	return false
}
