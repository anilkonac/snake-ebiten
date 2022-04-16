package game

import (
	"bytes"
	"math/rand"

	sound "github.com/anilkonac/snake-ebiten/game/resources/audio"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
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
	audioContext  *audio.Context
	slapPlayer    *audio.Player
	eatingPlayers []*audio.Player
	eatingSound   int // For debugging
	musicPlayer   *audio.Player
	playMusic     bool = true
)

func init() {
	prepareAudio()
}

func prepareAudio() {
	audioContext = audio.NewContext(sampleRate)
	eatingPlayers = make([]*audio.Player, 0)

	eatingPlayers = append(eatingPlayers, createPlayer(sound.Eating))
	eatingPlayers = append(eatingPlayers, createPlayer(sound.Eating2))
	eatingPlayers = append(eatingPlayers, createPlayer(sound.Eating3))
	eatingPlayers = append(eatingPlayers, createPlayer(sound.Eating4))

	slapPlayer = createPlayer(sound.Slap)

	musicPlayer = createMusicPlayer(sound.Music)
	musicPlayer.SetVolume(0.3)
	musicPlayer.Play()
}

func createPlayer(src []byte) *audio.Player {
	stream, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(src))
	if err != nil {
		panic(err)
	}

	player, err := audioContext.NewPlayer(stream)
	if err != nil {
		panic(err)
	}

	return player
}

func createMusicPlayer(src []byte) *audio.Player {
	stream, err := mp3.Decode(audioContext, bytes.NewReader(sound.Music))
	if err != nil {
		panic(err)
	}

	player, err := audioContext.NewPlayer(stream)
	if err != nil {
		panic(err)
	}

	return player
}

func play(sound soundE) {
	if sound >= soundTotal {
		panic("Wrong sound enum")
	}

	var player *audio.Player
	if sound == soundEating {
		randIndex := rand.Intn(len(eatingPlayers))
		player = eatingPlayers[randIndex]
		eatingSound = randIndex
		showSlap = false
	} else if sound == soundSlap {
		player = slapPlayer
		showSlap = true
	}

	player.Rewind()
	player.Play()
}
