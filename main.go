package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"landan-desktop-fyne/ui"
)

func main() {
	oApp := app.New()
	oWindow := oApp.NewWindow("Hello Fyne")
	oWindow.SetContent(ui.NewContent())
	oWindow.Resize(fyne.NewSize(1280, 720))
	oWindow.ShowAndRun()
}
