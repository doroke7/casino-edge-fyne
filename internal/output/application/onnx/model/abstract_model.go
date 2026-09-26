package outputApplicationOnnxModel

import (
	outputApplicationOnnx "landan-desktop-fyne/internal/output/application/onnx"
)

// AbstractModel 把 onnx 的共用資源（Context）包一層，
// model 這層目前不需要額外欄位。
type AbstractModel struct {
	*outputApplicationOnnx.AbstractOnnx
}

func NewAbstractModel(oAbstractOnnx *outputApplicationOnnx.AbstractOnnx) *AbstractModel {
	return &AbstractModel{
		AbstractOnnx: oAbstractOnnx,
	}
}
