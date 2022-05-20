package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	numSnakes   = 10
	minTurnTime = 3.0
	maxTurnTime = 10.0
)

type titleScene struct {
	snakes        []*snake
	sinceLastTurn map[int]float32
	time          float32
}

func newTitleScreen() *titleScene {
	scene := &titleScene{
		snakes:        make([]*snake, numSnakes),
		sinceLastTurn: make(map[int]float32),
	}

	for iSnake := 0; iSnake < numSnakes; iSnake++ {
		scene.snakes[iSnake] = newSnakeRandDirLoc()
	}
	return scene
}

func (t *titleScene) update() {
	t.time += deltaTime

	t.handleSettingsInputs()

	for iSnake, snake := range t.snakes {

		// Determine if turn
		if sinceLastTurn := t.sinceLastTurn[iSnake]; (sinceLastTurn >= minTurnTime) && (sinceLastTurn <= maxTurnTime) {

		}

		// Determine the new direction.
		dirCurrent := snake.lastDirection()
		dirNew := dirCurrent

		randResult := rand.Float32()
		if dirCurrent.isVertical() {
			if randResult < 0.1 {
				dirNew = directionLeft
			} else if randResult < 0.2 {
				dirNew = directionRight
			}
		} else {
			if randResult < 0.1 {
				dirNew = directionUp
			} else if randResult < 0.2 {
				dirNew = directionDown
			}
		}

		if dirNew != dirCurrent {
			// Create a new turn and take it
			newTurn := newTurn(dirCurrent, dirNew)
			snake.turnTo(newTurn, false)
		}

		snake.update(eatingAnimStartDistance)
	}
}

func (t *titleScene) draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)

	for _, snake := range t.snakes {
		// Draw the snake
		for unit := snake.unitHead; unit != nil; unit = unit.next {
			draw(screen, unit)
		}
	}

	drawFPS(screen)
}

func (t *titleScene) handleSettingsInputs() {
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		if musicState == musicOn {
			musicState = musicMuted
			playerMusic.Pause()
		} else if musicState == musicMuted {
			musicState = musicOn
			playerMusic.Play()
		}
	}
}
