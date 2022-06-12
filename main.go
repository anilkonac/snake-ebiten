package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	g "github.com/anilkonac/snake-ebiten/game"
	"github.com/hajimehoshi/ebiten/v2"
)

var cpuprofile = flag.String("cpuprofile", "", "write cput profile to 'file'")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	ebiten.SetWindowSize(g.ScreenSize())
	ebiten.SetWindowTitle("Ssnake")
	ebiten.RunGame(g.NewGame())
}

//go:generate file2byteslice -input game/shader/basic.go -output game/shader/basic_go.go -package shader -var Basic
//go:generate file2byteslice -input game/shader/round.go -output game/shader/round_go.go -package shader -var Round
//go:generate file2byteslice -input game/shader/score.go -output game/shader/score_go.go -package shader -var Score
//go:generate file2byteslice -input game/shader/snakehead.go -output game/shader/snakehead_go.go -package shader -var SnakeHead
//go:generate file2byteslice -input game/shader/title.go -output game/shader/title_go.go -package shader -var Title
