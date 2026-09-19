package main

import (
	"errors"
	"fmt"

	"github.com/ncruces/zenity"
)

type layoutDialogResult struct {
	path string
	save bool
	err  error
}

func (g *Game) openLayoutDialog(save bool) {
	if g.layoutDialog != nil {
		return
	}
	options := []zenity.Option{
		zenity.FileFilters{
			{Name: "JSON layouts", Patterns: []string{"*.json"}},
			{Name: "All files", Patterns: []string{"*"}},
		},
	}
	if g.layoutPath != "" {
		options = append(options, zenity.Filename(g.layoutPath))
	} else if save {
		options = append(options, zenity.Filename("layout.json"))
	}
	result := make(chan layoutDialogResult, 1)
	g.layoutDialog = result
	go func() {
		var path string
		var err error
		if save {
			path, err = zenity.SelectFileSave(append(options, zenity.Title("Save layout"), zenity.ConfirmOverwrite())...)
		} else {
			path, err = zenity.SelectFile(append(options, zenity.Title("Load layout"))...)
		}
		result <- layoutDialogResult{path: path, save: save, err: err}
	}()
}

func (g *Game) finishLayoutDialog() {
	select {
	case result := <-g.layoutDialog:
		g.layoutDialog = nil
		if errors.Is(result.err, zenity.ErrCanceled) {
			return
		}
		if result.err != nil {
			g.editorStatus = fmt.Sprintf("File dialog failed: %v", result.err)
			return
		}
		if result.path == "" {
			return
		}
		action := "saved"
		if result.save {
			if err := saveLayout(result.path, g.layout); err != nil {
				g.editorStatus = err.Error()
				return
			}
		} else {
			layout, err := loadLayout(result.path)
			if err != nil {
				g.editorStatus = err.Error()
				return
			}
			g.layout = layout
			g.currSelectedButton = -1
			action = "loaded"
		}
		g.layoutPath = result.path
		if err := saveLayout(g.lastLayoutPath, g.layout); err != nil {
			g.editorStatus = fmt.Sprintf("Layout %s, but could not remember it: %v", action, err)
			return
		}
		g.editorStatus = fmt.Sprintf("Layout %s.", action)
	default:
	}
}
