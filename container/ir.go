package container

import (
	inputApplicationIrInference "landan-desktop-fyne/internal/input/application/ir/inference"
)

type IRContainer struct {
	IrInferenceDie   *inputApplicationIrInference.DieHandler
	IrInferencePoker *inputApplicationIrInference.PokerHandler
	IrInferenceDisk  *inputApplicationIrInference.DiskHandler
}

func InitIRContainer() *IRContainer {
	return &IRContainer{
		IrInferenceDie:   inputApplicationIrInference.NewDieHandler(),
		IrInferencePoker: inputApplicationIrInference.NewPokerHandler(),
		IrInferenceDisk:  inputApplicationIrInference.NewDiskHandler(),
	}
}
