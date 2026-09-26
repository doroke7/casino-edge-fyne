package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	oApp := app.New()
	oWindow := oApp.NewWindow("Hello Fyne")

	oNameEntry := widget.NewEntry()
	oNameEntry.SetPlaceHolder("輸入你的名字")

	oGreetingLabel := widget.NewLabel("")

	oButton := widget.NewButton("打招呼", func() {
		oGreetingLabel.SetText("你好," + oNameEntry.Text + "!")
	})

	oWindow.SetContent(container.NewVBox(oNameEntry, oButton, oGreetingLabel))
	oWindow.Resize(fyne.NewSize(300, 150))
	oWindow.ShowAndRun()
}
