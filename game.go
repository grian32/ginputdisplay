package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

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

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBox(screen, 30, 0, 30, "W")
	g.drawBox(screen, 0, 30, 30, "A")
	g.drawBox(screen, 30, 30, 30, "S")
	g.drawBox(screen, 60, 30, 30, "D")

	g.drawBox(screen, 60, 0, 90, "GRAB")

	g.drawBox(screen, 90, 30, 50, "DASH")
	g.drawBox(screen, 140, 30, 50, "JUMP")
}

func (g *Game) drawBox(screen *ebiten.Image, x, y float32, width float32, msg string) {
	vector.StrokeRect(
		screen,
		x, y,
		width, 30,
		3,
		color.White,
		false,
	)

	strWidth, strHeight := text.Measure(msg, face, 1)
	textX := (float64(width) - strWidth) / 2.0
	textY := (float64(30) - strHeight) / 2.0
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x)+textX, float64(y)+textY)

	text.Draw(screen, msg, face, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
