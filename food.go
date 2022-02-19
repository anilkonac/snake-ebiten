package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	return newFood(float64(rand.Intn(screenWidth)), float64(rand.Intn(screenHeight)))
}

func (f food) draw(screen *ebiten.Image) {
	if !f.isActive {
		return
	}

	for _, rect := range f.rects {
		ebitenutil.DrawRect(screen, rect.x, rect.y, rect.width, rect.height, colorFood)
	}
}
