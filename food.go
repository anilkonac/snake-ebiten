package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type pixel uint16

const (
	foodLength     = 10
	foodLengthHalf = foodLength / 2.0
)

type food struct {
	rects []rectF64
}

func newFood(centerX, centerY pixel) *food {
	newFood := &food{
		rects: make([]rectF64, 0, 4),
	}

	pureRect := rectF64{
		x:      float64(centerX) - foodLengthHalf,
		y:      float64(centerY) - foodLengthHalf,
		width:  foodLength,
		height: foodLength,
	}
	pureRect.slice(&newFood.rects)

	return newFood
}

func newFoodRandLoc() *food {
	return newFood(pixel(rand.Intn(screenWidth)), pixel(rand.Intn(screenHeight)))
}

func (f food) draw(screen *ebiten.Image) {
	for _, rect := range f.rects {
		ebitenutil.DrawRect(screen, rect.x, rect.y, rect.width, rect.height, colorFood)
	}
}
