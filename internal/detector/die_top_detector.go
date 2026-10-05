package detector

import (
	"fmt"

	bootstrap "landan-desktop-fyne/bootstrap"
	domain "landan-desktop-fyne/internal/domain"
	pkgInference "landan-desktop-fyne/pkg/inference"
)

// DieTopDetector 用 onnx 跑 die 的偵測模型，onnx 的細節都在 AbstractDetector。
type DieTopDetector struct {
	*AbstractDetector
	threshold float32
}

// NewDieTopDetector 從 config/inference.yaml 的 detect.die.top 讀模型路徑與信心門檻。
func NewDieTopDetector(oInference *pkgInference.Inference) (*DieTopDetector, error) {
	if bootstrap.CONFIG.INFERENCE.DETECT.DIE.TOP.ONNX == "" {
		return nil, fmt.Errorf("inference.detect.die.top.onnx is empty (is config/inference.yaml filled in? run from the project root)")
	}

	oAbstractDetector, err := NewAbstractDetector(oInference, bootstrap.CONFIG.INFERENCE.DETECT.DIE.TOP.ONNX, bootstrap.CONFIG.INFERENCE.DETECT.DIE.TOP.OPENVINO)
	if err != nil {
		return nil, err
	}

	return &DieTopDetector{AbstractDetector: oAbstractDetector, threshold: bootstrap.CONFIG.INFERENCE.DETECT.DIE.TOP.THRESHOLD}, nil
}

func (oSelf *DieTopDetector) Recognize(aImage []byte) ([]*domain.Die, error) {
	aDetections, err := oSelf.AbstractDetector.Recognize(aImage, oSelf.threshold)
	if err != nil {
		return nil, err
	}

	var aDies []*domain.Die
	for _, oDetection := range aDetections {
		aDies = append(aDies, &domain.Die{
			X:          oDetection.X,
			Y:          oDetection.Y,
			Width:      oDetection.Width,
			Height:     oDetection.Height,
			Confidence: oDetection.Confidence,
		})
	}

	return aDies, nil
}

// main desktop compose up supervisor
// main desktop compose up
// main desktop
