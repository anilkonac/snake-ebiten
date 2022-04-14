package game

import (
	"bytes"

	sound "github.com/anilkonac/snake-ebiten/game/resources/audio"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type soundE uint8

const (
	soundEating soundE = iota
	soundSlap
	soundTotal
)
const sampleRate = 44100

var (
	audioContext *audio.Context
	audioPlayers map[soundE]*audio.Player
)

func prepareAudio() {
	audioContext = audio.NewContext(sampleRate)
	audioPlayers = make(map[soundE]*audio.Player, soundTotal)

	addPlayer(soundEating, sound.Eating)
	addPlayer(soundSlap, sound.Slap)
}

func addPlayer(sound soundE, src []byte) {
	stream, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(src))
	if err != nil {
		panic(err)
	}
	audioPlayers[sound], err = audioContext.NewPlayer(stream)
	if err != nil {
		panic(err)
	}
}

func play(sound soundE) {
	if sound >= soundTotal {
		panic("Wrong sound enum")
	}

	player := audioPlayers[sound]
	player.Rewind()
	player.Play()
}
