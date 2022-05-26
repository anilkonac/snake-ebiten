package tools

import (
	"math"
)

type VecI struct {
	X, Y int
}

func (v VecI) To32() Vec32 {
	return Vec32{float32(v.X), float32(v.Y)}
}

func (v VecI) To64() Vec64 {
	return Vec64{float64(v.X), float64(v.Y)}
}

type Vec32 struct {
	X, Y float32
}

func (v Vec32) To64() Vec64 {
	return Vec64{float64(v.X), float64(v.Y)}
}

type Vec64 struct {
	X, Y float64
}

func (v Vec64) To32() Vec32 {
	return Vec32{float32(v.X), float32(v.Y)}
}

func (v Vec64) Floor() Vec64 {
	return Vec64{math.Floor(v.X), math.Floor(v.Y)}
}

func Distance(a, b Vec64) float64 {
	distX := a.X - b.X
	distY := a.Y - b.Y
	return math.Hypot(distX, distY)
}
