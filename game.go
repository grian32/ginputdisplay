package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"sync"

	ebitenbackend "github.com/AllenDang/cimgui-go/backend/ebiten-backend"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	mu     sync.RWMutex
	keys   map[evdev.EvCode]bool
	dev    *evdev.Device
	layout DisplayLayout
	ui     *ebitenbackend.EbitenBackend
	showUI bool
	clicks int
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
	g.ui = ebitenbackend.NewEbitenBackend()
	g.ui.CreateWindow("ginputdisplay", 640, 480)
	imgui.CurrentIO().SetIniFilename("")
	g.showUI = true
}

func (g *Game) Destroy() {
	g.dev.Close()
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.showUI = !g.showUI
	}
	g.ui.BeginFrame()
	if g.showUI {
		imgui.SetNextWindowPosV(imgui.NewVec2(20, 120), imgui.CondOnce, imgui.NewVec2(0, 0))
		imgui.SetNextWindowSizeV(imgui.NewVec2(280, 100), imgui.CondOnce)
		if imgui.Begin("ImGui") {
			imgui.Text("F1: show/hide this window")
			if imgui.Button("Click me") {
				g.clicks++
			}
			imgui.SameLine()
			imgui.Text(fmt.Sprintf("Clicks: %d", g.clicks))
		}
		imgui.End()
	}
	g.ui.EndFrame()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.mu.RLock()
	for _, box := range g.layout.Boxes {
		g.drawBox(screen, box, g.keys[box.code])
	}
	g.mu.RUnlock()
	g.ui.Draw(screen)
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

func (g *Game) drawBox(screen *ebiten.Image, box Box, fill bool) {
	x, y := float32(box.X), float32(box.Y)
	width, height := float32(box.Width), float32(box.Height)
	if fill {
		vector.FillRect(
			screen,
			x, y,
			width, height,
			color.White,
			false,
		)
	} else {
		vector.StrokeRect(
			screen,
			x, y,
			width, height,
			3,
			color.White,
			false,
		)
	}

	strWidth, strHeight := text.Measure(box.Label, face, 1)
	textX := (float64(width) - strWidth) / 2.0
	textY := (float64(height) - strHeight) / 2.0
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x)+textX, float64(y)+textY)
	op.ColorScale.Scale(0, 104, 163, 255)

	text.Draw(screen, box.Label, face, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ui.Layout(g.layout.Width, g.layout.Height)
}
