package game

import (
	"math"
)

type vecI struct {
	x, y int
}

func (v vecI) to32() vec32 {
	return vec32{float32(v.x), float32(v.y)}
}

type vec32 struct {
	x, y float32
}

func (v vec32) to64() vec64 {
	return vec64{float64(v.x), float64(v.y)}
}

type vec64 struct {
	x, y float64
}

func (v vec64) to32() vec32 {
	return vec32{float32(v.x), float32(v.y)}
}

func (v vec64) floor() vec64 {
	return vec64{math.Floor(v.x), math.Floor(v.y)}
}
