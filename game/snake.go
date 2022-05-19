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
)

// Snake parameters
const (
	snakeHeadCenterX        = halfScreenWidth
	snakeHeadCenterY        = halfScreenHeight
	snakeSpeedInitial       = 275
	snakeSpeedFinal         = 250
	snakeLength             = 240
	snakeWidth              = 30
	eatingAnimStartDistance = 120
	radiusSnake             = snakeWidth / 2.0
	radiusMouth             = radiusSnake * 0.625
)

// Snake collision tolerances must be an integer or false collisions will occur.
const (
	toleranceDefault    = 2 //snakeWidth / 16.0
	toleranceScreenEdge = radiusSnake
)

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

func newSnake(headCenter vec64, direction directionT) *snake {
	if direction >= directionTotal {
		panic("direction parameter is invalid.")
	}
	if headCenter.x > ScreenWidth {
		panic("Initial x position of the snake is off-screen.")
	}
	if headCenter.y > ScreenHeight {
		panic("Initial y position of the snake is off-screen.")
	}
	if isVertical := direction.isVertical(); (isVertical && (snakeLength > ScreenHeight)) ||
		(!isVertical && (snakeLength > ScreenWidth)) {
		panic("Initial snake intersects itself.")
	}

	initialUnit := newUnit(headCenter, snakeLength, direction, &colorSnake1)

	snake := &snake{
		speed:    snakeSpeedInitial,
		unitHead: initialUnit,
		unitTail: initialUnit,
	}

	return snake
}

func newSnakeRandDir(headCenter vec64) *snake {
	direction := directionT(rand.Intn(int(directionTotal)))
	return newSnake(headCenter, direction)
}

func (s *snake) update(distToFood float32) {
	moveDistance := s.speed * deltaTime

	// if the snake has moved a safe distance after the last turn, take the next turn in the queue.
	if (len(s.turnQueue) > 0) && (s.distAfterTurn+toleranceDefault >= snakeWidth) {
		var nextTurn *turn
		nextTurn, s.turnQueue = s.turnQueue[0], s.turnQueue[1:] // Pop front
		s.turnTo(nextTurn, true)
	}

	s.updateHead(moveDistance, distToFood)
	s.updateTail(moveDistance, distToFood)
}

func (s *snake) updateHead(dist float64, distToFood float32) {
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
		s.unitHead.update(distToFood)
	}

	s.distAfterTurn += dist
}

func (s *snake) updateTail(dist float64, distToFood float32) {
	decreaseAmount := dist
	if s.growthRemaining > 0 {
		// Calculate the tail reduction with the square function so that the growth doesn't look ugly.
		// f(x) = 0.8 + 3.0*(x-0.5)^2)/4.0 where  0 <= x <= 1
		growthCompletion := 1 - s.growthRemaining/s.growthTarget
		decreaseAmount *= (0.8 + 3.0*(growthCompletion-0.5)*(growthCompletion-0.5)/4.0)
		s.growthRemaining -= (dist - decreaseAmount)
	} else {
		s.growthTarget = s.growthRemaining
	}

	// Decrease tail length
	s.unitTail.length -= decreaseAmount

	// Delete tail if its length is less than width of the snake
	if (s.unitTail.prev != nil) && (s.unitTail.length <= snakeWidth) {
		s.unitTail.prev.length += s.unitTail.length
		s.unitTail = s.unitTail.prev
		s.unitTail.next = nil
	}

	s.unitTail.update(distToFood)
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
	newHead := newUnit(oldHead.headCenter, 0, newTurn.directionTo, newColor)

	// Add the new head unit to the beginning of the unit doubly linked list.
	newHead.next = oldHead
	oldHead.prev = newHead
	s.unitHead = newHead

	// Update prev turn
	s.turnPrev = newTurn
}

func (s *snake) checkIntersection(intersected *bool) {
	*intersected = false
	curUnit := s.unitHead.next
	if curUnit == nil {
		return
	}

	var tolerance float32 = toleranceDefault
	if len(curUnit.rectsCollision) > 1 { // If second unit is on an edge
		tolerance = toleranceScreenEdge // To avoid false collisions on screen edges
	}

	for curUnit != nil {
		if collides(s.unitHead, curUnit, tolerance) {
			*intersected = true
			playSoundHit()
			return
		}
		curUnit = curUnit.next
	}
}

func (s *snake) grow() {
	// Compute the new growth and add to the remaining growth value.
	// f(x)=50+5*log2(x/10.0+1)
	newGrowth := 50.0 + 5.0*math.Log2(float64(s.foodEaten)/10.0+1.0)
	s.growthRemaining += newGrowth
	s.growthTarget += newGrowth
	s.foodEaten++

	// Update snake speed
	// f(x)=250+25/e^(0.0075x)
	s.speed = snakeSpeedFinal + (snakeSpeedInitial-snakeSpeedFinal)/math.Exp(0.0075*float64(s.foodEaten))
}

func (s *snake) lastDirection() directionT {
	// if the turn queue is not empty, return the direction of the last turn to be taken.
	if queueLength := len(s.turnQueue); queueLength > 0 {
		return s.turnQueue[queueLength-1].directionTo
	}

	// return current head direction
	return s.unitHead.direction
}
