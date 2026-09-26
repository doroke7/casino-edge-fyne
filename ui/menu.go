package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func NewMainMenu(oWindow fyne.Window) *fyne.MainMenu {

	// "Settings…" is a special label: on macOS Fyne moves it into the app menu.
	oSettingsItem := fyne.NewMenuItem("Settings…", func() {
		ShowSettings()
	})
	oSettingsItem.Shortcut = &desktop.CustomShortcut{
		KeyName:  fyne.KeyComma,
		Modifier: fyne.KeyModifierSuper,
	}

	oCloseItem := fyne.NewMenuItem("關閉視窗", func() {
		oWindow.Close()
	})
	oCloseItem.Shortcut = &desktop.CustomShortcut{
		KeyName:  fyne.KeyW,
		Modifier: fyne.KeyModifierSuper,
	}

	oFileMenu := fyne.NewMenu("檔案", oSettingsItem, fyne.NewMenuItemSeparator(), oCloseItem)

	oEditMenu := fyne.NewMenu("編輯",
		newEditItem(oWindow, "剪下", &fyne.ShortcutCut{Clipboard: fyne.CurrentApp().Clipboard()}),
		newEditItem(oWindow, "複製", &fyne.ShortcutCopy{Clipboard: fyne.CurrentApp().Clipboard()}),
		newEditItem(oWindow, "貼上", &fyne.ShortcutPaste{Clipboard: fyne.CurrentApp().Clipboard()}),
		fyne.NewMenuItemSeparator(),
		newEditItem(oWindow, "全選", &fyne.ShortcutSelectAll{}),
	)

	oCameraItem := fyne.NewMenuItem("開啟攝影機", ToggleCamera)
	oSettingsMenu := fyne.NewMenu("設定", oCameraItem)

	oMainMenu := fyne.NewMainMenu(oFileMenu, oEditMenu, oSettingsMenu)

	SetCameraStateHandler(func(bOn bool) {
		if bOn {
			oCameraItem.Label = "關閉攝影機"
		} else {
			oCameraItem.Label = "開啟攝影機"
		}
		oMainMenu.Refresh()
	})

	return oMainMenu
}

// newEditItem sends a clipboard/selection shortcut to the currently focused widget.
func newEditItem(oWindow fyne.Window, sLabel string, oShortcut fyne.Shortcut) *fyne.MenuItem {

	oItem := fyne.NewMenuItem(sLabel, func() {
		if oTarget, ok := oWindow.Canvas().Focused().(fyne.Shortcutable); ok {
			oTarget.TypedShortcut(oShortcut)
		}
	})
	oItem.Shortcut = oShortcut

	return oItem
}
