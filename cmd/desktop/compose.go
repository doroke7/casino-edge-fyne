// Package desktop 放宿主機(desktop)端的命令,目前只有 docker compose 呼叫的入口 compose:
//
//	desktop compose --project-name=NAME up SERVICE
//	desktop compose --project-name=NAME down SERVICE
//	desktop compose metadata
package desktop

import (
	"github.com/spf13/cobra"

	"landan-desktop-fyne/cmd/desktop/down"
	"landan-desktop-fyne/cmd/desktop/metadata"
	"landan-desktop-fyne/cmd/desktop/up"
)

var Command = &cobra.Command{
	Use:   "compose",
	Short: "docker compose 呼叫的入口",
}

func init() {
	// docker compose 會帶 --project-name
	Command.PersistentFlags().String("project-name", "", "docker compose 專案名稱")

	Command.AddCommand(metadata.Command, up.Command, down.Command)
}
