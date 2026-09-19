package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	layoutPath := flag.String("layout", "", "layout JSON file (defaults to the last loaded layout)")
	flag.Parse()
	lastPath, err := lastLayoutPath()
	if err != nil {
		log.Fatal(err)
	}
	path := *layoutPath
	if path == "" {
		path = lastPath
	}
	layout, err := loadLayout(path)
	status := ""
	if err != nil && *layoutPath == "" {
		if !errors.Is(err, os.ErrNotExist) {
			status = fmt.Sprintf("Could not restore layout: %v", err)
			log.Print(status)
		}
		layout, err = loadLayout("")
	}
	if err != nil {
		log.Fatal(err)
	}
	if *layoutPath != "" {
		if err := saveLayout(lastPath, layout); err != nil {
			status = fmt.Sprintf("Layout loaded, but could not remember it: %v", err)
			log.Print(status)
		}
	}
	game := &Game{layout: layout, layoutPath: *layoutPath, lastLayoutPath: lastPath, editorStatus: status}
	game.Setup()
	go game.readKeyboard()
	defer game.Destroy()
	if err := ebiten.RunGameWithOptions(game, &ebiten.RunGameOptions{ScreenTransparent: true}); err != nil {
		log.Fatal(err)
	}
}
