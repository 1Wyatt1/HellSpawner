package hspaletteeditor

import (
	"path/filepath"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
)

type PaletteEditor struct {
	hseditor.Editor
	palette d2interface.Palette
	path    string
}

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	palette, err := d2dat.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &PaletteEditor{
		path:    filepath.Base(pathEntry.FullPath),
		palette: palette,
	}

	return result, nil
}

func (e *PaletteEditor) GetWindowTitle() string {
	return e.path + "##" + e.GetId()
}

func (e *PaletteEditor) Render() {
	if !e.Visible {
		return
	}

	if e.ToFront {
		e.ToFront = false
		imgui.SetNextWindowFocus()
	}

	g.Window(e.GetWindowTitle()).IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Pos(360, 30).Layout(g.Layout{
		hswidget.PaletteGrid(e.GetId()+"_grid", e.palette.GetColors()),
	})
}

func (e *PaletteEditor) Cleanup() {
	e.Visible = false
}
