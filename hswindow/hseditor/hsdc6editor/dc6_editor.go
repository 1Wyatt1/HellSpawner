package hsdc6editor

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	dc6, err := d2dc6.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &DC6Editor{
		dc6: dc6,
	}

	result.Path = pathEntry

	return result, nil
}

type DC6Editor struct {
	hseditor.Editor
	dc6 *d2dc6.DC6
}

func (e *DC6Editor) Render() {
	if !e.Visible {
		return
	}

	if e.ToFront {
		e.ToFront = false
		imgui.SetNextWindowFocus()
	}

	g.Window(e.GetWindowTitle()).IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		hswidget.DC6Viewer(e.Path.GetUniqueId(), e.dc6),
		g.Custom(func() {
			e.Focused = imgui.IsWindowFocused(0)
		}),
	})
}

func (e *DC6Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DC6 Editor").Layout(g.Layout{
		g.MenuItem("Add to project").OnClick(func() {}),
		g.MenuItem("Remove from project").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Import from file...").OnClick(func() {}),
		g.MenuItem("Export to file...").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Close").OnClick(func() {
			e.Visible = false
		}),
	})

	*l = append(*l, m)
}
