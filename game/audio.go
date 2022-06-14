/*
Copyright (C) 2022 Anıl Konaç

This file is part of snake-ebiten.

snake-ebiten is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

snake-ebiten is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with snake-ebiten. If not, see <https://www.gnu.org/licenses/>.
*/

package game

import (
	"bytes"
	"math/rand"
	"time"

	c "github.com/anilkonac/snake-ebiten/game/core"
	sound "github.com/anilkonac/snake-ebiten/resources/audio"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	sampleRate    = 44100
	volumeEating  = 0.45
	volumeMusic   = 0.4
	volumeHit     = 1.0
	probEatingA   = 0.74
	musicCheckSec = 2
)

type stateMusic uint8

const (
	musicOn stateMusic = iota
	musicPaused
	musicMuted
)

var (
	audioContext  audio.Context
	playerHit     audio.Player
	playerMusic   audio.Player
	playerEatingA audio.Player
	playerEatingB audio.Player
	musicState    stateMusic
	playSounds    = true
)

func init() {
	prepareAudio()
	playerMusic.Play()
	go repeatMusic()
}

func prepareAudio() {
	audioContext = *audio.NewContext(sampleRate)

	playerEatingA = *createPlayer(sound.Eating, volumeEating)
	playerEatingB = *createPlayer(sound.Eating2, volumeEating)

	playerHit = *createPlayer(sound.Hit, volumeHit)

	playerMusic = *createMusicPlayer(sound.Music)
	playerMusic.SetVolume(volumeMusic)
}

func createPlayer(src []byte, volume float64) *audio.Player {
	stream, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(src))
	c.Panic(err)

	player, err := audioContext.NewPlayer(stream)
	c.Panic(err)

	player.SetVolume(volume)
	return player
}

func createMusicPlayer(src []byte) *audio.Player {
	stream, err := vorbis.DecodeWithSampleRate(sampleRate, bytes.NewReader(src))
	c.Panic(err)

	player, err := audioContext.NewPlayer(stream)
	c.Panic(err)

	return player
}

func playSoundEating() {
	if !playSounds {
		return
	}

	var player *audio.Player
	if rand.Float32() < probEatingA {
		player = &playerEatingA
	} else {
		player = &playerEatingB
	}
	player.Rewind()
	player.Play()
}

func playSoundHit() {
	if !playSounds {
		return
	}
	playerHit.Rewind()
	playerHit.Play()
}

// Goroutine
func repeatMusic() {
	ticker := time.NewTicker(time.Second * musicCheckSec)
	for range ticker.C {
		if (musicState == musicOn) && !playerMusic.IsPlaying() {
			playerMusic.Rewind()
			playerMusic.Play()
		}
	}
}
