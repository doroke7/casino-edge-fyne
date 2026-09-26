package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

func SetupTray(oApp fyne.App) {

	oDesktopApp, ok := oApp.(desktop.App)
	if !ok {
		return
	}

	oSettingsItem := fyne.NewMenuItem("設定", func() {
		fyne.Do(ShowSettings)
	})
	oSettingsItem.Icon = theme.SettingsIcon()

	oDesktopApp.SetSystemTrayMenu(fyne.NewMenu("", oSettingsItem))
	oDesktopApp.SetSystemTrayIcon(theme.SettingsIcon())
}
