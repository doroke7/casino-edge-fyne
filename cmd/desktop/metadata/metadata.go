package metadata

import (
	"fmt"

	"github.com/spf13/cobra"

	"landan-desktop-fyne/internal/logger"
)

var Command = &cobra.Command{
	Use:   "metadata",
	Short: "輸出參數說明給 docker compose",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), logger.Metadata)
	},
}
