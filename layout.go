package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mikegio27/go-evdev"
)

//go:embed assets/default-layout.json
var defaultLayout []byte

func lastLayoutPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("find config directory: %w", err)
	}
	return filepath.Join(dir, "ginputdisplay", "last-layout.json"), nil
}

func saveLayout(path string, layout DisplayLayout) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("enter a layout file path")
	}
	data, err := json.MarshalIndent(layout, "", "  ")
	if err != nil {
		return fmt.Errorf("encode layout: %w", err)
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create layout directory: %w", err)
	}
	file, err := os.CreateTemp(dir, ".layout-*.json")
	if err != nil {
		return fmt.Errorf("save layout %q: %w", path, err)
	}
	defer os.Remove(file.Name())
	_, writeErr := file.Write(append(data, '\n'))
	closeErr := file.Close()
	if writeErr != nil {
		return fmt.Errorf("write layout %q: %w", path, writeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close layout %q: %w", path, closeErr)
	}
	if err := os.Rename(file.Name(), path); err != nil {
		return fmt.Errorf("save layout %q: %w", path, err)
	}
	return nil
}

type DisplayLayout struct {
	Width  int   `json:"width"`
	Height int   `json:"height"`
	Boxes  []Box `json:"boxes"`
}

type Box struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Label  string `json:"label"`
	Key    string `json:"key"`
	code   evdev.EvCode
}

func loadLayout(path string) (DisplayLayout, error) {
	data := defaultLayout
	if path != "" {
		var err error
		data, err = os.ReadFile(path)
		if err != nil {
			return DisplayLayout{}, fmt.Errorf("read layout %q: %w", path, err)
		}
	}
	var layout DisplayLayout
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&layout); err != nil {
		return layout, fmt.Errorf("decode layout: %w", err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return layout, fmt.Errorf("layout must contain exactly one JSON object")
	}
	if layout.Width <= 0 || layout.Height <= 0 {
		return layout, fmt.Errorf("layout width and height must be positive")
	}
	for i := range layout.Boxes {
		box := &layout.Boxes[i]
		if box.Width <= 0 || box.Height <= 0 {
			return layout, fmt.Errorf("box %d (%q): width and height must be positive", i+1, box.Label)
		}
		code, ok := evdev.EvCodeByName(box.Key)
		if !ok || (!strings.HasPrefix(box.Key, "KEY_") && !strings.HasPrefix(box.Key, "BTN_")) {
			return layout, fmt.Errorf("box %d (%q): unknown key %q; use an evdev name such as KEY_W", i+1, box.Label, box.Key)
		}
		box.code = code
	}
	return layout, nil
}
