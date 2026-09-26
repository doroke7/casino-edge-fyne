package up

import (
	"fmt"

	"github.com/spf13/cobra"

	"landan-desktop-fyne/bootstrap/launcher"
	"landan-desktop-fyne/internal/logger"
)

var Command = &cobra.Command{
	Use:   "up SERVICE",
	Short: "在背景啟動桌面程式",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nPid, bStarted, err := launcher.Start()
		if err != nil {
			logger.Error(args[0] + ": " + err.Error())
			return err
		}
		if !bStarted {
			logger.Info(fmt.Sprintf("%s: 已經在執行 (pid %d)", args[0], nPid))
			return nil
		}
		logger.Info(fmt.Sprintf("%s: 已啟動 (pid %d),日誌 %s", args[0], nPid, launcher.LogPath()))
		return nil
	},
}
