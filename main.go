package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("ginputdisplay")
	game := &Game{}
	game.Setup()
	go game.readKeyboard()
	defer game.Destroy()
	if err := ebiten.RunGameWithOptions(game, &ebiten.RunGameOptions{ScreenTransparent: true}); err != nil {
		log.Fatal(err)
		game.Destroy()
	}
}
