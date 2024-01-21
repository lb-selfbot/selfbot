//go:build lbmobile

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	App    fyne.App
	Window fyne.Window
}

func New(title string) *GUI {
	a := app.New()

	g := &GUI{
		App:    a,
		Window: a.NewWindow(title),
	}

	g.SetupGUI()

	return g
}

func (g *GUI) SetupGUI() {
	toolbar := widget.NewMenu(fyne.NewMenu("Menu test", fyne.NewMenuItem("item1", func() {}), fyne.NewMenuItem("item1", func() {}), fyne.NewMenuItem("item1", func() {})))

	content := container.NewBorder(toolbar, nil, nil, nil, widget.NewLabel("Content"))

	g.Window.SetContent(
		container.NewVBox(
			widget.NewLabel("Hello Fyne!"),
			content,
			widget.NewButton("Quit", func() {
				g.App.Quit()
			}),
		),
	)
}

func (g *GUI) Main() {
	g.Window.ShowAndRun()
}
