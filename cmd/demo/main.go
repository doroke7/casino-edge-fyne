package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello Fyne")

	name := widget.NewEntry()
	name.SetPlaceHolder("輸入你的名字")

	greeting := widget.NewLabel("")

	button := widget.NewButton("打招呼", func() {
		greeting.SetText("你好," + name.Text + "!")
	})

	w.SetContent(container.NewVBox(name, button, greeting))
	w.Resize(fyne.NewSize(300, 150))
	w.ShowAndRun()
}
