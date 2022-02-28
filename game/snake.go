/*
snake-ebiten
Copyright (C) 2022 Anıl Konaç

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package game

import (
	"math/rand"
	"time"
)

type directionT uint8
type snakeLengthT uint16

const (
	safeDist            = 0.5
	toleranceDefault    = 0.001
	toleranceScreenEdge = halfSnakeWidth
)

const (
	directionUp directionT = iota
	directionDown
	directionLeft
	directionRight
	directionTotal
)

func (d directionT) isVertical() bool {
	if d >= directionTotal {
		panic("wrong direction")
	}
	return (d == directionUp) || (d == directionDown)
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
	if centerX > ScreenWidth {
		panic("Initial x position of the snake is off-screen.")
	}
	if centerY > ScreenHeight {
		panic("Initial y position of the snake is off-screen.")
	}
	if isVertical := direction.isVertical(); (isVertical && (snakeLength > ScreenHeight)) ||
		(!isVertical && (snakeLength > ScreenWidth)) {
		panic("Initial snake intersects itself.")
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
	direction := directionT(rand.Intn(int(directionTotal)))
	return newSnake(centerX, centerY, direction, speed, snakeLength)
}

func (s *snake) update() {
	moveDistance := float64(s.speed) * deltaTime

	// if the snake has moved a safe distance after the last turn, take the next turn in the queue.
	if (len(s.turnQueue) > 0) && (s.distAfterTurn-snakeWidth >= safeDist) {
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

	if (s.unitHead != s.unitTail) || (s.remainingGrowth > 0) { // Avoid unnecessary updates
		s.unitHead.creteRects() // Update rectangles of this unit
	}

	s.distAfterTurn += dist
}

func (s *snake) updateTail(dist float64) {
	if s.remainingGrowth > 0 {
		s.remainingGrowth -= dist
		return
	}

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

	s.unitTail.creteRects() // Update rectangles of this unit
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
	curUnit := s.unitHead.next
	if curUnit == nil {
		return false
	}

	tolerance := toleranceDefault
	if len(curUnit.rects) > 1 { // If second unit is on an edge
		tolerance = toleranceScreenEdge // To avoid false collisions on screen edges
	}

	for curUnit != nil {
		if collides(s.unitHead, curUnit, tolerance) {
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
