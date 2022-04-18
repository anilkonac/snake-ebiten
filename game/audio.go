/*
snake-ebiten
Copyright (C) 2022 Anıl Konaç

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
	soundPickup soundE = iota
	soundHit
	soundTotal
)
const sampleRate = 44100

var (
	audioContext  *audio.Context
	playerHit     *audio.Player
	playersEating [2]*audio.Player
	playerMusic   *audio.Player
	playMusic     bool = true
)

func init() {
	prepareAudio()
}

func prepareAudio() {
	audioContext = audio.NewContext(sampleRate)
	// pickupPlayers = make([]*audio.Player, 0)

	playersEating[0] = createPlayer(sound.Eating, 0.6)
	playersEating[1] = createPlayer(sound.Eating2, 0.55)

	playerHit = createPlayer(sound.Hit, 1.0)

	playerMusic = createMusicPlayer(sound.Music)
	playerMusic.SetVolume(0.5)
	playerMusic.Play()
}

func createPlayer(src []byte, volume float64) *audio.Player {
	stream, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(src))
	if err != nil {
		panic(err)
	}

	player, err := audioContext.NewPlayer(stream)
	if err != nil {
		panic(err)
	}

	player.SetVolume(volume)
	return player
}

func createMusicPlayer(src []byte) *audio.Player {
	stream, err := mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(src))
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
	if sound == soundPickup {
		randIndex := rand.Intn(11)
		var eatingIndex uint8
		if randIndex < 8 {
			eatingIndex = 0
		} else {
			eatingIndex = 1
		}
		player = playersEating[eatingIndex]
		showSlap = false
	} else if sound == soundHit {
		player = playerHit
		showSlap = true
	}

	player.Rewind()
	player.Play()
}
