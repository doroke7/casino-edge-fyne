package inputApplicationIrInference

import (
	pbIrTableInference "landan-desktop-fyne/pb/ir/table/inference"
)

type PokerHandler struct {
	pbIrTableInference.UnimplementedPokerServiceServer
}

func NewPokerHandler() *PokerHandler {
	return &PokerHandler{}
}
