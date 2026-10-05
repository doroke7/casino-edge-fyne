package detector

import (
	"fmt"
	"sort"

	bootstrap "landan-desktop-fyne/bootstrap"
	domain "landan-desktop-fyne/internal/domain"
	pkgInference "landan-desktop-fyne/pkg/inference"
)

// PokerCardDetector 用 onnx 跑撲克牌的偵測模型（只有一個類別 Card），onnx 的細節都在 AbstractDetector。
type PokerCardDetector struct {
	*AbstractDetector
}

// NewPokerCardDetector 從 config/inference.yaml 的 detect.poker.card 讀模型路徑與信心門檻。
func NewPokerCardDetector(oInference *pkgInference.Inference) (*PokerCardDetector, error) {
	if bootstrap.CONFIG.INFERENCE.DETECT.POKER.CARD.ONNX == "" {
		return nil, fmt.Errorf("inference.detect.poker.card.onnx is empty (is config/inference.yaml filled in? run from the project root)")
	}

	oAbstractDetector, err := NewAbstractDetector(oInference, bootstrap.CONFIG.INFERENCE.DETECT.POKER.CARD.ONNX, bootstrap.CONFIG.INFERENCE.DETECT.POKER.CARD.OPENVINO)
	if err != nil {
		return nil, err
	}

	return &PokerCardDetector{AbstractDetector: oAbstractDetector}, nil
}

// Recognize 回傳圖中的撲克牌，依框的中心 x 由小到大排序。
func (oSelf *PokerCardDetector) Recognize(aImage []byte) ([]*domain.Poker, error) {
	aDetections, err := oSelf.AbstractDetector.Recognize(aImage, bootstrap.CONFIG.INFERENCE.DETECT.POKER.CARD.THRESHOLD)
	if err != nil {
		return nil, err
	}

	var aPokers []*domain.Poker
	for _, oDetection := range aDetections {
		aPokers = append(aPokers, &domain.Poker{
			X:          oDetection.X,
			Y:          oDetection.Y,
			Width:      oDetection.Width,
			Height:     oDetection.Height,
			Confidence: oDetection.Confidence,
		})
	}
	sort.SliceStable(aPokers, func(iLeft, iRight int) bool {
		return aPokers[iLeft].X*2+aPokers[iLeft].Width < aPokers[iRight].X*2+aPokers[iRight].Width
	})

	return aPokers, nil
}
