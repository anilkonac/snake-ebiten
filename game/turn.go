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

type turn struct {
	directionTo   directionT
	isTurningLeft bool
}

func newTurn(directionFrom, directionTo directionT) *turn {
	newTurn := &turn{
		directionTo: directionTo,
	}

	// Determine the direction of rotation.
	if (directionFrom == directionUp && directionTo == directionLeft) ||
		(directionFrom == directionLeft && directionTo == directionDown) ||
		(directionFrom == directionDown && directionTo == directionRight) ||
		(directionFrom == directionRight && directionTo == directionUp) {
		newTurn.isTurningLeft = true
	}

	return newTurn
}
