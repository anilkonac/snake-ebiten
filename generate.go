package main

//go:generate file2byteslice -input game/shaders/basic.go -output game/shaders/basic_go.go -package shaders -var Basic
//go:generate file2byteslice -input game/shaders/hollow.go -output game/shaders/hollow_go.go -package shaders -var Hollow
//go:generate file2byteslice -input game/shaders/round.go -output game/shaders/round_go.go -package shaders -var Round
