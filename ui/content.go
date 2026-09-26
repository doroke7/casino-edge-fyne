package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewContent() *fyne.Container {

	oNameEntry := widget.NewEntry()
	oNameEntry.SetPlaceHolder("輸入你的名字")

	oGreetingLabel := widget.NewLabel("")

	oButton := widget.NewButton("打招呼", func() {
		oGreetingLabel.SetText("你好," + oNameEntry.Text + "!")
	})

	return container.NewVBox(oNameEntry, oButton, oGreetingLabel)
}