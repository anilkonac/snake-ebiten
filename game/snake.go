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
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type directionT uint8
type snakeLengthT uint16

const (
	toleranceDefault    = snakeWidth / 16.0
	toleranceScreenEdge = halfSnakeWidth
	toleranceFood       = snakeWidth / 4.0
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
	speed           float64
	unitHead        *unit
	unitTail        *unit
	turnPrev        *turn
	turnQueue       []*turn
	distAfterTurn   float64
	growthRemaining float64
	growthTarget    float64
	foodEaten       uint8
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newSnake(centerX, centerY float64, direction directionT, snakeLength snakeLengthT) *snake {
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
		speed:    snakeSpeedInitial,
		unitHead: initialUnit,
		unitTail: initialUnit,
	}

	return snake
}

func newSnakeRandDir(centerX, centerY float64, snakeLength snakeLengthT) *snake {
	direction := directionT(rand.Intn(int(directionTotal)))
	return newSnake(centerX, centerY, direction, snakeLength)
}

func (s *snake) update() {
	moveDistance := s.speed * deltaTime

	// if the snake has moved a safe distance after the last turn, take the next turn in the queue.
	if (len(s.turnQueue) > 0) && (s.distAfterTurn+toleranceDefault >= snakeWidth) {
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

	if s.unitHead != s.unitTail { // Avoid unnecessary updates
		s.unitHead.creteRects() // Update rectangles of this unit
	}

	s.distAfterTurn += dist
}

func (s *snake) updateTail(dist float64) {
	decreaseAmount := dist
	if s.growthRemaining > 0 {
		// Calculate the tail reduction with the square function so that the growth doesn't look ugly.
		// f(x) = 0.75 + (x-0.5)^2 where  0 <= x <= 1
		growthCompletion := 1 - s.growthRemaining/s.growthTarget
		decreaseAmount *= (0.75 + (growthCompletion-0.5)*(growthCompletion-0.5))
		s.growthRemaining -= decreaseAmount
	} else {
		s.growthTarget = s.growthRemaining
	}

	// Decrease tail length
	s.unitTail.length -= decreaseAmount

	// Delete tail if its length is less than width of the snake
	if (s.unitTail.prev != nil) && (s.unitTail.length <= snakeWidth) {
		s.unitTail = s.unitTail.prev
		s.unitTail.length += snakeWidth
		s.unitTail.next = nil
	}

	s.unitTail.creteRects() // Update rectangles of this unit
}

func (s *snake) draw(dst *ebiten.Image) {
	// Create units between head centers
	var drawableUnits []*unit

	for unit := s.unitHead; unit.next != nil; unit = unit.next {
		newUnit := newUnit(unit.headCenterX, unit.headCenterY, unit.length+snakeWidth, unit.direction, unit.color)
		drawableUnits = append(drawableUnits, newUnit)
	}
	drawableUnits = append(drawableUnits, s.unitTail)

	// Draw these units
	for iUnit := len(drawableUnits) - 1; iUnit >= 0; iUnit-- {
		draw(dst, drawableUnits[iUnit])
	}
}

func (s *snake) turnTo(newTurn *turn, isFromQueue bool) {
	if !isFromQueue {
		// Check if the new turn is dangerous (twice same turns rapidly).
		if (s.turnPrev != nil) &&
			(s.turnPrev.isTurningLeft == newTurn.isTurningLeft) &&
			(s.distAfterTurn+toleranceDefault <= snakeWidth) {
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

	var tolerance float32 = toleranceDefault
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

	// Compute the new growth and add to the remaining growth value.
	// f(x)=20/(e^(0.0125x))
	increasePercent := 20.0 / math.Exp(0.0125*float64(s.foodEaten))
	curGrowth := totalLength * increasePercent / 100.0
	s.growthRemaining += curGrowth
	s.growthTarget += curGrowth
	s.foodEaten++

	// Update snake speed
	// f(x)=275+25/e^(0.010625x)
	s.speed = snakeSpeedFinal + (snakeSpeedInitial-snakeSpeedFinal)/math.Exp(0.010625*float64(s.foodEaten))
}

func (s *snake) lastDirection() directionT {
	// if the turn queue is not empty, return the direction of the last turn to be taken.
	if queueLength := len(s.turnQueue); queueLength > 0 {
		return s.turnQueue[queueLength-1].directionTo
	}

	// return current head direction
	return s.unitHead.direction
}
