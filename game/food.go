package game

import (
	"image/color"
	"math/rand"
)

const (
	foodLength     = snakeWidth / 2.0
	halfFoodLength = halfSnakeWidth
)

type food struct {
	isActive bool
	rects    []rectF64
}

func newFood(centerX, centerY float64) *food {
	newFood := &food{
		rects: make([]rectF64, 0, 4),
	}

	// Create a rectangle to use in drawing and eating logic.
	pureRect := rectF64{
		x:      centerX - halfFoodLength,
		y:      centerY - halfFoodLength,
		width:  foodLength,
		height: foodLength,
	}
	// Split this rectangle if it is on a screen edge.
	pureRect.split(&newFood.rects)

	return newFood
}

func newFoodRandLoc() *food {
	return newFood(float64(rand.Intn(ScreenWidth)), float64(rand.Intn(ScreenHeight)))
}

// Implement slicer interface
// --------------------------
func (f food) slice() []rectF64 {
	return f.rects
}

// Implement collidable interface
// ------------------------------
func (f food) collEnabled() bool {
	return true
}

// Implement drawable interface
// ------------------------------
func (f food) drawEnabled() bool {
	return f.isActive
}

func (f food) Color() color.Color {
	return colorFood
}
