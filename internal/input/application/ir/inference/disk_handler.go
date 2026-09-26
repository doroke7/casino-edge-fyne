package inputApplicationIrInference

import (
	pbIrTableInference "landan-desktop-fyne/pb/ir/table/inference"
)

type DiskHandler struct {
	pbIrTableInference.UnimplementedDiskServiceServer
}

func NewDiskHandler() *DiskHandler {
	return &DiskHandler{}
}
