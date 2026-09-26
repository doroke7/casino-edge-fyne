// Package imagerecognition is the `desktop image-recognition` command: a headless gRPC server.
package imagerecognition

import (
	"fmt"
	"net"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"landan-desktop-fyne/internal/helper"
	"landan-desktop-fyne/internal/imagerecognition"
	pb "landan-desktop-fyne/proto/ir"
)

var Command = &cobra.Command{
	Use:   "image-recognition",
	Short: "啟動影像辨識 gRPC 服務(不開視窗)",
	RunE: func(cmd *cobra.Command, args []string) error {
		sAddr, err := cmd.Flags().GetString("addr")
		if err != nil {
			return err
		}
		oListener, err := net.Listen("tcp", sAddr)
		if err != nil {
			return err
		}

		oServer := grpc.NewServer()
		pb.RegisterImageRecognitionServiceServer(oServer, imagerecognition.NewServer(imagerecognition.BasicRecognizer{}))

		oCtx, fnStop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer fnStop()
		go func() {
			<-oCtx.Done()
			oServer.GracefulStop()
		}()

		helper.Info(fmt.Sprintf("image-recognition gRPC 服務監聽 %s", oListener.Addr()))
		return oServer.Serve(oListener)
	},
}

func init() {
	Command.Flags().String("addr", "127.0.0.1:50051", "gRPC 監聽位址")
}
