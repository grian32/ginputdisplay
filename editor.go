package main

import (
	"fmt"
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/mikegio27/go-evdev"
)

var keyOptions = []evdev.EvCode{
	evdev.KEY_A, evdev.KEY_B, evdev.KEY_C, evdev.KEY_D, evdev.KEY_E,
	evdev.KEY_F, evdev.KEY_G, evdev.KEY_H, evdev.KEY_I, evdev.KEY_J,
	evdev.KEY_K, evdev.KEY_L, evdev.KEY_M, evdev.KEY_N, evdev.KEY_O,
	evdev.KEY_P, evdev.KEY_Q, evdev.KEY_R, evdev.KEY_S, evdev.KEY_T,
	evdev.KEY_U, evdev.KEY_V, evdev.KEY_W, evdev.KEY_X, evdev.KEY_Y, evdev.KEY_Z,
	evdev.KEY_LEFTSHIFT, evdev.KEY_RIGHTSHIFT,
	evdev.KEY_LEFTCTRL, evdev.KEY_RIGHTCTRL,
	evdev.KEY_LEFTALT, evdev.KEY_RIGHTALT,
	evdev.KEY_LEFTMETA, evdev.KEY_RIGHTMETA,
	evdev.KEY_KP0, evdev.KEY_KP1, evdev.KEY_KP2, evdev.KEY_KP3, evdev.KEY_KP4,
	evdev.KEY_KP5, evdev.KEY_KP6, evdev.KEY_KP7, evdev.KEY_KP8, evdev.KEY_KP9,
	evdev.KEY_KPDOT, evdev.KEY_KPPLUS, evdev.KEY_KPMINUS,
	evdev.KEY_KPASTERISK, evdev.KEY_KPSLASH, evdev.KEY_KPENTER,
	evdev.KEY_KPEQUAL, evdev.KEY_NUMLOCK,
}

func drawKeyCombo(box *Box) {
	if !imgui.BeginCombo("Key", strings.TrimPrefix(box.Key, "KEY_")) {
		return
	}
	for _, code := range keyOptions {
		name := evdev.CodeName(evdev.EV_KEY, code)
		selected := box.code == code
		if imgui.SelectableBoolV(strings.TrimPrefix(name, "KEY_"), selected, 0, imgui.NewVec2(0, 0)) {
			box.Key = name
			box.code = code
		}
		if selected {
			imgui.SetItemDefaultFocus()
		}
	}
	imgui.EndCombo()
}

func drawSizingModifiers(box *Box) {
	imgui.Text(fmt.Sprintf("X: %d, Y: %d, Width: %d, Height: %d", box.X, box.Y, box.Width, box.Height))
	imgui.PushItemFlag(imgui.ItemFlagsButtonRepeat, true)
	if imgui.SmallButton("-##x") {
		box.X--
	}
	imgui.SameLine()
	if imgui.SmallButton("+##x") {
		box.X++
	}
	imgui.SameLine()
	if imgui.SmallButton("-##y") {
		box.Y--
	}
	imgui.SameLine()
	if imgui.SmallButton("+##y") {
		box.Y++
	}
	imgui.SameLine()
	if imgui.SmallButton("-##width") && box.Width > 1 {
		box.Width--
	}
	imgui.SameLine()
	if imgui.SmallButton("+##width") {
		box.Width++
	}
	imgui.SameLineV(185.0, -0.0)
	if imgui.SmallButton("-##height") && box.Height > 1 {
		box.Height--
	}
	imgui.SameLine()
	if imgui.SmallButton("+##height") {
		box.Height++
	}
	imgui.PopItemFlag()
}

func drawLabelInput(box *Box, index int) {
	imgui.PushIDInt(int32(index))
	imgui.InputTextWithHint("Label", "", &box.Label, 0, nil)
	imgui.PopID()
}

func (g *Game) drawLayoutControls() {
	if imgui.Button("Add new button") {
		g.layout.Boxes = append(g.layout.Boxes, Box{
			X: 40, Y: 40, Width: 30, Height: 30,
			Label: "A", Key: "KEY_A", code: evdev.KEY_A,
		})
		g.currSelectedButton = len(g.layout.Boxes) - 1
	}
	imgui.SameLine()
	imgui.BeginDisabledV(g.currSelectedButton == -1)
	if imgui.Button("Remove selected") {
		i := g.currSelectedButton
		g.layout.Boxes = append(g.layout.Boxes[:i], g.layout.Boxes[i+1:]...)
		g.currSelectedButton = -1
	}
	imgui.EndDisabled()

	imgui.BeginDisabledV(g.layoutDialog != nil)
	if imgui.Button("Save layout...") {
		g.openLayoutDialog(true)
	}
	imgui.SameLine()
	if imgui.Button("Load layout...") {
		g.openLayoutDialog(false)
	}
	imgui.EndDisabled()

	if g.editorStatus != "" {
		imgui.TextWrapped(g.editorStatus)
	}
	imgui.Separator()
}

func (g *Game) drawEditor() {
	imgui.SetNextWindowPosV(imgui.NewVec2(20, 120), imgui.CondOnce, imgui.NewVec2(0, 0))
	imgui.SetNextWindowSizeV(imgui.NewVec2(360, 270), imgui.CondOnce)
	if imgui.Begin("Editor") {
		imgui.Text("F1: show/hide this window")
		g.drawLayoutControls()
		if g.currSelectedButton != -1 {
			currButton := &g.layout.Boxes[g.currSelectedButton]
			drawSizingModifiers(currButton)
			drawKeyCombo(currButton)
			drawLabelInput(currButton, g.currSelectedButton)
		}
	}
	imgui.End()
}
