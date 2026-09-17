package main

import (
	"bytes"
	"image/color"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mikegio27/go-evdev"

	_ "embed"
)

//go:embed assets/GrindyBrush.otf
var fontData []byte

var fontSource *text.GoTextFaceSource

var face *text.GoTextFace

func init() {
	var err error
	fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatal(err)
	}

	face = &text.GoTextFace{
		Source: fontSource,
		Size:   14,
	}
}

type Game struct {
	mu   sync.RWMutex
	keys map[evdev.EvCode]bool
	dev  *evdev.Device
}

func (g *Game) Setup() {
	g.keys = make(map[evdev.EvCode]bool)
	dev, err := evdev.Open(
		"/dev/input/by-id/usb-Logitech_LogiG_MKeyboard-event-kbd",
	)
	if err != nil {
		log.Fatal(err)
	}
	g.dev = dev
}

func (g *Game) Destroy() {
	g.dev.Close()
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.mu.Lock()
	g.drawBox(screen, 30, 0, 30, "W", g.keys[evdev.KEY_W])
	g.drawBox(screen, 0, 30, 30, "A", g.keys[evdev.KEY_A])
	g.drawBox(screen, 30, 30, 30, "S", g.keys[evdev.KEY_S])
	g.drawBox(screen, 60, 30, 30, "D", g.keys[evdev.KEY_D])

	g.drawBox(screen, 60, 0, 90, "GRAB", g.keys[evdev.KEY_LEFTSHIFT])

	g.drawBox(screen, 90, 30, 50, "DASH", g.keys[evdev.KEY_KP2])
	g.drawBox(screen, 140, 30, 50, "JUMP", g.keys[evdev.KEY_KP3])
	g.mu.Unlock()
}

func (g *Game) readKeyboard() {
	for {
		ev, err := g.dev.ReadOne()
		if err != nil {
			log.Println(err)
			return
		}

		if ev.Type != evdev.EV_KEY {
			continue
		}

		g.mu.Lock()

		switch ev.Value {
		case 0:
			g.keys[ev.Code] = false
		case 1:
			g.keys[ev.Code] = true
		case 2:
			g.keys[ev.Code] = true
		}

		g.mu.Unlock()
	}
}

func (g *Game) drawBox(screen *ebiten.Image, xGiven, yGiven float32, width float32, msg string, fill bool) {
	x := xGiven + 40
	y := yGiven + 40
	if fill {
		vector.FillRect(
			screen,
			x, y,
			width, 30,
			color.White,
			false,
		)
	} else {
		vector.StrokeRect(
			screen,
			x, y,
			width, 30,
			3,
			color.White,
			false,
		)
	}

	strWidth, strHeight := text.Measure(msg, face, 1)
	textX := (float64(width) - strWidth) / 2.0
	textY := (float64(30) - strHeight) / 2.0
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x)+textX, float64(y)+textY)
	op.ColorScale.Scale(0, 104, 163, 255)

	text.Draw(screen, msg, face, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
