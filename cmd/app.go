package cmd

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/spf13/cobra"

	"landan-desktop-fyne/internal/bootstrap"
	"landan-desktop-fyne/ui"
)

var rootCmd = &cobra.Command{
	Use:           "app",
	Short:         "Landan desktop app built with Fyne",
	SilenceUsage:  true, // a config error is not a usage error
	SilenceErrors: true, // Execute prints it once
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := bootstrap.CONFIG.Validate(); err != nil {
			return err
		}

		oApp := app.New()
		ui.SetupTray(oApp)
		oApp.Lifecycle().SetOnStopped(ui.ShutdownCamera)
		oWindow := oApp.NewWindow("Hello Fyne")
		oWindow.SetMainMenu(ui.NewMainMenu(oWindow))
		oWindow.SetContent(ui.NewContent())
		oWindow.Resize(fyne.NewSize(1280, 720))
		oWindow.ShowAndRun()
		return nil
	},
}

// Execute runs the root command; called from the top-level main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
