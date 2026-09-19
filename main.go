package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	layoutPath := flag.String("layout", "", "layout JSON file (defaults to the embedded layout)")
	flag.Parse()
	layout, err := loadLayout(*layoutPath)
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("ginputdisplay")
	game := &Game{layout: layout}
	game.Setup()
	go game.readKeyboard()
	defer game.Destroy()
	if err := ebiten.RunGameWithOptions(game, &ebiten.RunGameOptions{ScreenTransparent: true}); err != nil {
		log.Fatal(err)
		game.Destroy()
	}
}
