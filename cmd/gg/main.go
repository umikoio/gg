package main

import (
	"fmt"

	"github.com/Kaamkiya/gg/internal/app/maze"
	
	"github.com/charmbracelet/huh"
)

func main() {
	var game string

	fmt.Println("gg - a tui for small offline games")

	huh.NewSelect[string]().
		Title("choose a game:").
		Options(
			huh.NewOption("maze", "maze"),
			huh.NewOption("2048", "2048"),
			//huh.NewOption("tetris", "tetris"),
		).
		Value(&game).
		Run()
	
	switch game {
	case "maze":
		maze.Run()
	default:
		panic("This game either doesn't exist or hasn't been implemented.")
	}
}
