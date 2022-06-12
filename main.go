package main

import (
	"fmt"
	"log"
	"os"

	g "github.com/anilkonac/snake-ebiten/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(g.ScreenSize())
	ebiten.SetWindowTitle("Ssnake")
	ebiten.RunGame(g.NewGame())

	f, err := os.OpenFile("results.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fTicks := float64(g.NumTicks)
	avgTPS := g.TPSSum / fTicks
	avgFPS := g.FPSSum / fTicks
	if _, err := f.WriteString(fmt.Sprintf("%d\t\t%f\t%f\t%f\t\t%f\t%f\t%f\n", g.NumTicks, g.TPSMin, g.TPSMax, avgTPS, g.FPSMin, g.FPSMax, avgFPS)); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("TicksElapsed: %d\n", g.NumTicks)
	fmt.Printf("TPS | Min: %.2f, Max: %.2f, Avg: %.2f\n", g.TPSMin, g.TPSMax, avgTPS)
	fmt.Printf("FPS | Min: %.2f, Max: %.2f, Avg: %.2f\n", g.FPSMin, g.FPSMax, avgFPS)
}

//go:generate file2byteslice -input game/shader/basic.go -output game/shader/basic_go.go -package shader -var Basic
//go:generate file2byteslice -input game/shader/round.go -output game/shader/round_go.go -package shader -var Round
//go:generate file2byteslice -input game/shader/score.go -output game/shader/score_go.go -package shader -var Score
//go:generate file2byteslice -input game/shader/snakehead.go -output game/shader/snakehead_go.go -package shader -var SnakeHead
//go:generate file2byteslice -input game/shader/title.go -output game/shader/title_go.go -package shader -var Title
