package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var oSettingsWindow fyne.Window

func ShowSettings() {

	if oSettingsWindow != nil {
		oSettingsWindow.RequestFocus()
		return
	}

	oSettingsWindow = fyne.CurrentApp().NewWindow("設定")
	oSettingsWindow.SetContent(widget.NewLabel("設定內容"))
	oSettingsWindow.Resize(fyne.NewSize(480, 320))
	oSettingsWindow.SetOnClosed(func() {
		oSettingsWindow = nil
	})
	oSettingsWindow.Show()
}
