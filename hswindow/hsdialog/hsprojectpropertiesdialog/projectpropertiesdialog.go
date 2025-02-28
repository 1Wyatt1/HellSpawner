package hsprojectpropertiesdialog

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/OpenDiablo2/HellSpawner/hscommon"

	"github.com/OpenDiablo2/HellSpawner/hsconfig"

	"github.com/OpenDiablo2/HellSpawner/hswidget"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog"
)

const (
	removeItemButtonPath = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_delete.png"
	upItemButtonPath     = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_up.png"
	downItemButtonPath   = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_down.png"
)

type ProjectPropertiesDialog struct {
	hsdialog.Dialog

	removeIconTexture          *g.Texture
	upIconTexture              *g.Texture
	downIconTexture            *g.Texture
	project                    hsproject.Project
	config                     *hsconfig.Config
	onProjectPropertiesChanged func(project hsproject.Project)
	auxMPQs, auxMPQNames       []string

	mpqSelectDlgIndex      int
	mpqSelectDialogVisible bool
}

func Create(onProjectPropertiesChanged func(project hsproject.Project)) *ProjectPropertiesDialog {
	result := &ProjectPropertiesDialog{
		onProjectPropertiesChanged: onProjectPropertiesChanged,
		mpqSelectDialogVisible:     false,
	}

	hscommon.CreateTextureFromFileAsync(removeItemButtonPath, func(texture *g.Texture) {
		result.removeIconTexture = texture
	})

	hscommon.CreateTextureFromFileAsync(upItemButtonPath, func(texture *g.Texture) {
		result.upIconTexture = texture
	})

	hscommon.CreateTextureFromFileAsync(downItemButtonPath, func(texture *g.Texture) {
		result.downIconTexture = texture
	})

	return result
}

func (p *ProjectPropertiesDialog) Show(project *hsproject.Project, config *hsconfig.Config) {
	p.config = config
	p.project = *project
	p.auxMPQs = config.GetAuxMPQs()
	p.auxMPQNames = make([]string, len(p.auxMPQs))

	for idx := range p.auxMPQNames {
		p.auxMPQNames[idx] = filepath.Base(p.auxMPQs[idx])
	}
	p.Dialog.Show()
}

func (p *ProjectPropertiesDialog) Render() {
	if !p.Visible {
		return
	}

	canSave := len(strings.TrimSpace(p.project.ProjectName)) > 0

	hswidget.ModalDialog("Select Auxiliary MPQ##ProjectPropertiesSelectAuxMPQDialog", &p.mpqSelectDialogVisible, g.Layout{
		g.Child("ProjectPropertiesSelectAuxMPQDialogLayout").Size(300, 200).Layout(g.Layout{
			g.ListBox("ProjectPropertiesSelectAuxMPQDialogItems", p.auxMPQNames).Border(false).OnChange(func(selectedIndex int) {
				p.mpqSelectDlgIndex = selectedIndex
			}).OnDClick(func(selectedIndex int) {
				p.mpqSelectDialogVisible = false
				p.addAuxMpq(p.auxMPQs[selectedIndex])
			}),
		}),
		g.Line(
			g.Button("Add Selected...##ProjectPropertiesSelectAuxMPQDialogAddSelected").OnClick(func() {
				p.addAuxMpq(p.auxMPQs[p.mpqSelectDlgIndex])
				p.mpqSelectDialogVisible = false
			}),
			g.Button("Cancel##ProjectPropertiesSelectAuxMPQDialogCancel").OnClick(func() {
				p.mpqSelectDialogVisible = false
			}),
		),
	})

	if !p.mpqSelectDialogVisible {
		hswidget.ModalDialog("Project Properties##ProjectPropertiesDialog", &p.Visible, g.Layout{
			g.Line(
				g.Child("ProjectPropertiesLayout").Size(300, 250).Layout(g.Layout{
					g.Label("Project Name:"),
					g.InputText("##ProjectPropertiesDialogProjectName", &p.project.ProjectName).Size(250),
					g.Label("Description:"),
					g.InputText("##ProjectPropertiesDialogDescription", &p.project.Description).Size(250),
					g.Label("Author:"),
					g.InputText("##ProjectPropertiesDialogAuthor", &p.project.Author).Size(250),
				}),
				g.Child("ProjectPropertiesLayout2").Size(300, 250).Layout(g.Layout{
					g.Label("Auxiliary MPQs:"),
					g.Child("ProjectPropertiesAuxMpqLayoutGroup").Border(false).Size(0, 180).Layout(g.Layout{
						g.Custom(func() {
							imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{})
							imgui.PushStyleColor(imgui.StyleColorBorder, imgui.Vec4{})
							imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{})
							for idx := range p.project.AuxiliaryMPQs {
								if idx >= len(p.project.AuxiliaryMPQs) {
									break
								}
								g.Line(
									g.Custom(func() { imgui.PushID(fmt.Sprintf("ProjectPropertiesAddAuxMpqRemove_%d", idx)) }),
									g.ImageButton(p.removeIconTexture).Size(16, 16).OnClick(func() {
										copy(p.project.AuxiliaryMPQs[idx:], p.project.AuxiliaryMPQs[idx+1:])
										p.project.AuxiliaryMPQs = p.project.AuxiliaryMPQs[:len(p.project.AuxiliaryMPQs)-1]
									}),
									g.Custom(func() {
										imgui.PopID()
										imgui.PushID(fmt.Sprintf("ProjectPropertiesAddAuxMpqDown_%d", idx))
									}),
									g.ImageButton(p.downIconTexture).Size(16, 16).OnClick(func() {
										if idx < len(p.project.AuxiliaryMPQs)-1 {
											p.project.AuxiliaryMPQs[idx], p.project.AuxiliaryMPQs[idx+1] = p.project.AuxiliaryMPQs[idx+1], p.project.AuxiliaryMPQs[idx]
										}
									}),
									g.Custom(func() {
										imgui.PopID()
										imgui.PushID(fmt.Sprintf("ProjectPropertiesAddAuxMpqUp_%d", idx))
									}),
									g.ImageButton(p.upIconTexture).Size(16, 16).OnClick(func() {
										if idx > 0 {
											p.project.AuxiliaryMPQs[idx-1], p.project.AuxiliaryMPQs[idx] = p.project.AuxiliaryMPQs[idx], p.project.AuxiliaryMPQs[idx-1]
										}
									}),
									g.Custom(func() { imgui.PopID() }),
									g.Dummy(8, 0),
									g.Label(p.project.AuxiliaryMPQs[idx]),
								).Build()
							}
							imgui.PopStyleVar()
							imgui.PopStyleColorV(2)
						}),
					}),
					g.Button("Add Auxiliary MPQ...##ProjectPropertiesAddAuxMpq").OnClick(p.onAddAuxMpqClicked),
				}),
			),
			g.Line(
				g.Custom(func() {
					if !canSave {
						imgui.PushStyleVarFloat(imgui.StyleVarAlpha, 0.5)
					}
				}),
				g.Button("Save##ProjectPropertiesDialogSave").OnClick(p.onSaveClicked),
				g.Custom(func() {
					if !canSave {
						imgui.PopStyleVar()
					}
				}),
				g.Button("Cancel##ProjectPropertiesDialogCancel").OnClick(p.onCancelClicked),
			),
		},
		)
	}
}

func (p *ProjectPropertiesDialog) onSaveClicked() {
	if len(strings.TrimSpace(p.project.ProjectName)) == 0 {
		return
	}

	p.onProjectPropertiesChanged(p.project)
	p.Visible = false
}

func (p *ProjectPropertiesDialog) onCancelClicked() {
	p.Visible = false
}

func (p *ProjectPropertiesDialog) onAddAuxMpqClicked() {
	p.mpqSelectDialogVisible = true
}

func (p *ProjectPropertiesDialog) addAuxMpq(mpqPath string) {
	relPath, err := filepath.Rel(p.config.AuxiliaryMpqPath, mpqPath)
	if err != nil {
		log.Fatal(err)
	}

	for idx := range p.project.AuxiliaryMPQs {
		if p.project.AuxiliaryMPQs[idx] == relPath {
			return
		}
	}

	p.project.AuxiliaryMPQs = append(p.project.AuxiliaryMPQs, relPath)
}
