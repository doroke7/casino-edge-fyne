package cmd

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/spf13/cobra"

	"landan-desktop-fyne/ui"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Landan desktop app built with Fyne",
	Run: func(cmd *cobra.Command, args []string) {
		oApp := app.New()
		oWindow := oApp.NewWindow("Hello Fyne")
		oWindow.SetContent(ui.NewContent())
		oWindow.Resize(fyne.NewSize(1280, 720))
		oWindow.ShowAndRun()
	},
}

// Execute runs the root command; called from the top-level main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

