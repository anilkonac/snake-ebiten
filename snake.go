package main

const (
	directionUp uint8 = iota
	directionDown
	directionLeft
	directionRight
	directionTotal
)

type unit struct {
	posX      float32
	posY      float32
	direction uint8
	frontUnit *unit
	backUnit  *unit
}

type snake struct {
	speed uint8
	head  *unit
}
