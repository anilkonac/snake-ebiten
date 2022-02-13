package main

type turn struct {
	directionTo   directionT
	isTurningLeft bool
}

func newTurn(directionFrom, directionTo directionT) *turn {
	newTurn := &turn{
		directionTo: directionTo,
	}

	// Identify turning direction
	if (directionFrom == directionUp && directionTo == directionLeft) ||
		(directionFrom == directionLeft && directionTo == directionDown) ||
		(directionFrom == directionDown && directionTo == directionRight) ||
		(directionFrom == directionRight && directionTo == directionUp) {
		newTurn.isTurningLeft = true
	}

	return newTurn
}