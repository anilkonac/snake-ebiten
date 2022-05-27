package param

// Snake parameters
const (
	SnakeSpeedInitial       = 275
	SnakeSpeedFinal         = 250
	SnakeLength             = 240
	SnakeWidth              = 30
	EatingAnimStartDistance = 120
	RadiusSnake             = SnakeWidth / 2.0
	RadiusMouth             = RadiusSnake * 0.625
)

// Snake collision tolerances must be an integer or false collisions will occur.
const (
	ToleranceDefault    = 2 //param.SnakeWidth / 16.0
	ToleranceScreenEdge = RadiusSnake
)
