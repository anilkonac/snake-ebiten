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

package snake

type DirectionT uint8

const (
	DirectionUp DirectionT = iota
	DirectionDown
	DirectionLeft
	DirectionRight
	DirectionTotal
)

func (d DirectionT) IsVertical() bool {
	if d >= DirectionTotal {
		panic("wrong direction")
	}
	return (d == DirectionUp) || (d == DirectionDown)
}

type Turn struct {
	directionTo   DirectionT
	isTurningLeft bool
}

func NewTurn(directionFrom, directionTo DirectionT) *Turn {
	newTurn := &Turn{
		directionTo: directionTo,
	}

	// Determine the direction of rotation.
	if (directionFrom == DirectionUp && directionTo == DirectionLeft) ||
		(directionFrom == DirectionLeft && directionTo == DirectionDown) ||
		(directionFrom == DirectionDown && directionTo == DirectionRight) ||
		(directionFrom == DirectionRight && directionTo == DirectionUp) {
		newTurn.isTurningLeft = true
	}

	return newTurn
}
