package main

var numTurn uint = 0

type turn struct {
	no            uint
	directionFrom directionT
	directionTo   directionT
	isTurningLeft bool
}

func newTurn(directionFrom, directionTo directionT) *turn {
	newTurn := &turn{
		directionFrom: directionFrom,
		directionTo:   directionTo,
		no:            numTurn,
	}
	numTurn++

	// Determine turning direction
	if (directionFrom == directionUp && directionTo == directionLeft) ||
		(directionFrom == directionLeft && directionTo == directionDown) ||
		(directionFrom == directionDown && directionTo == directionRight) ||
		(directionFrom == directionRight && directionTo == directionUp) {
		newTurn.isTurningLeft = true
	}

	return newTurn
}
