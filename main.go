package main

import (
	"fmt"

	g "github.com/anilkonac/snake-ebiten/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(g.ScreenSize())
	ebiten.SetWindowTitle("Ssnake")
	ebiten.RunGame(g.NewGame())

	fmt.Printf("TicksElapsed: %d, Mean FPS: %.2f", g.NumTicks, g.FPSSum/float64(g.NumTicks))
}

//go:generate file2byteslice -input game/shader/basic.go -output game/shader/basic_go.go -package shader -var Basic
//go:generate file2byteslice -input game/shader/round.go -output game/shader/round_go.go -package shader -var Round
//go:generate file2byteslice -input game/shader/score.go -output game/shader/score_go.go -package shader -var Score
//go:generate file2byteslice -input game/shader/snakehead.go -output game/shader/snakehead_go.go -package shader -var SnakeHead
//go:generate file2byteslice -input game/shader/title.go -output game/shader/title_go.go -package shader -var Title
