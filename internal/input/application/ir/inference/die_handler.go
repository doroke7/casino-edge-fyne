package inputApplicationIrInference

import (
	pbIrTableInference "landan-desktop-fyne/pb/ir/table/inference"
)

type DieHandler struct {
	pbIrTableInference.UnimplementedDieServiceServer
}

func NewDieHandler() *DieHandler {
	return &DieHandler{}
}
