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
	g.Window.SetContent(
		container.NewVBox(
			widget.NewLabel("Hello Fyne!"),
			widget.NewButton("Quit", func() {
				g.App.Quit()
			}),
		),
	)
}

func (g *GUI) Main() {
	g.Window.ShowAndRun()
}
