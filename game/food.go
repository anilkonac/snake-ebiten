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
	active bool
	rects  []rectF64
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

// Implement collidable interface
// ------------------------------
func (f food) isActive() bool {
	return f.active
}

func (f food) rectSlice() []rectF64 {
	return f.rects
}

// Implement drawable interface
// ------------------------------
func (f food) Color() color.Color {
	return colorFood
}
