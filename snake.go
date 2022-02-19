package main

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type directionT uint8
type snakeLengthT uint16

const epsilonSafeDist = 0.5

const (
	directionUp directionT = iota
	directionDown
	directionLeft
	directionRight
	directionTotal
)

func isDirectionVertical(dir directionT) bool {
	if dir >= directionTotal {
		panic("wrong direction")
	}
	return (dir == directionUp) || (dir == directionDown)
}

type snake struct {
	speed           uint8
	unitHead        *unit
	unitTail        *unit
	turnPrev        *turn
	turnQueue       []*turn
	distAfterTurn   float64
	remainingGrowth float64
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

	initialUnit := newUnit(centerX, centerY, float64(snakeLength), direction, &colorSnake1)

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

	// if the snake has moved a safe distance after the last turn, take the next turn in the queue.
	if (len(s.turnQueue) > 0) && (s.distAfterTurn-snakeWidth >= epsilonSafeDist) {
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
		s.unitHead.creteRects() // Update rectangles of this unit
	}

	s.distAfterTurn += dist
}

func (s *snake) updateTail(dist float64) {
	if s.remainingGrowth > 0 {
		s.remainingGrowth -= dist
	} else {
		// Decrease tail length
		s.unitTail.length -= dist

		// Rotate tail if its length is less than width of the snake
		if (s.unitTail.prev != nil) && (s.unitTail.length <= snakeWidth) {
			s.unitTail.direction = s.unitTail.prev.direction
		}

		// Destroy tail unit if its length is not positive
		if s.unitTail.length <= 0 {
			s.unitTail = s.unitTail.prev
			s.unitTail.next = nil
		}
	}

	s.unitTail.creteRects() // Update rectangles of this unit
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
		// If there are turns in the queue then add the new turn to the queue as well.
		if len(s.turnQueue) > 0 {
			s.turnQueue = append(s.turnQueue, newTurn)
			return
		}
	}
	s.distAfterTurn = 0

	oldHead := s.unitHead

	// Decide on the color of the new head unit.
	newColor := &colorSnake1
	if debugUnits && (oldHead.color == &colorSnake1) {
		newColor = &colorSnake2
	}

	// Create a new head unit.
	newHead := newUnit(oldHead.headCenterX, oldHead.headCenterY, 0, newTurn.directionTo, newColor)

	// Add the new head unit to the beginning of the unit doubly linked list.
	newHead.next = oldHead
	oldHead.prev = newHead
	s.unitHead = newHead

	// Update prev turn
	s.turnPrev = newTurn
}

func (s *snake) checkIntersection() bool {
	if s.unitHead.next == nil {
		return false
	}

	// Skip the next unit of the head, since it is not possible for the head to intersect with it.
	curUnit := s.unitHead.next.next
	for curUnit != nil {
		if s.unitHead.intersects(curUnit) {
			return true
		}
		curUnit = curUnit.next
	}
	return false
}

func (s *snake) grow() {
	// Compute the total length of the snake.
	var totalLength float64
	for unit := s.unitHead; unit != nil; unit = unit.next {
		totalLength += unit.length
	}

	s.remainingGrowth += totalLength * lengthIncreasePercent / 100.0
}

func (s *snake) lastDirection() directionT {
	// if the turn queue is not empty, return the direction of the last turn to be taken.
	if queueLength := len(s.turnQueue); queueLength > 0 {
		return s.turnQueue[queueLength-1].directionTo
	}

	// return current head direction
	return s.unitHead.direction
}
