package registerIR

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	container "landan-desktop-fyne/container"

	pbIrTableInference "landan-desktop-fyne/pb/ir/table/inference"
)

func Init(oContainer *container.IRContainer) *grpc.Server {

	oGrpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    1 * time.Second,
				Timeout: 5 * time.Second,
			},
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             10 * time.Second,
				PermitWithoutStream: true,
			},
		),
	)
	pbIrTableInference.RegisterDieServiceServer(oGrpcServer, oContainer.IrInferenceDie)
	pbIrTableInference.RegisterPokerServiceServer(oGrpcServer, oContainer.IrInferencePoker)
	pbIrTableInference.RegisterDiskServiceServer(oGrpcServer, oContainer.IrInferenceDisk)

	return oGrpcServer
}
