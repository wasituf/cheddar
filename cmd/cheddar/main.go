package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wasituf/cheddar/internal/app"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Cheddar")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(app.New()); err != nil {
		log.Fatal(err)
	}
}
