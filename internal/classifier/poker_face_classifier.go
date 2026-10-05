package classifier

import (
	"fmt"

	bootstrap "landan-desktop-fyne/bootstrap"
	domain "landan-desktop-fyne/internal/domain"
	pkgInference "landan-desktop-fyne/pkg/inference"
)

// 撲克牌的面，順序見 pk-studio-ir-model 的 cfg/classify/poker/card-data.yaml。
var aPokerFaceNames = []string{"Front", "Flow", "Back"}

// PokerFaceClassifier 用 onnx 跑 撲克牌哪一面朝上（Front / Flow / Back），onnx 的細節都在 AbstractClassifier。
type PokerFaceClassifier struct {
	*AbstractClassifier
}

// NewPokerFaceClassifier 從 config/inference.yaml 的 classify.poker.card 讀模型路徑。
func NewPokerFaceClassifier(oInference *pkgInference.Inference) (*PokerFaceClassifier, error) {
	if bootstrap.CONFIG.INFERENCE.CLASSIFY.POKER.CARD.ONNX == "" {
		return nil, fmt.Errorf("inference.classify.poker.card.onnx is empty (is config/inference.yaml filled in? run from the project root)")
	}

	oAbstractClassifier, err := NewAbstractClassifier(oInference, bootstrap.CONFIG.INFERENCE.CLASSIFY.POKER.CARD.ONNX, bootstrap.CONFIG.INFERENCE.CLASSIFY.POKER.CARD.OPENVINO)
	if err != nil {
		return nil, err
	}

	return &PokerFaceClassifier{AbstractClassifier: oAbstractClassifier}, nil
}

func (oSelf *PokerFaceClassifier) Classify(aImage []byte) (*domain.PokerFace, error) {
	sName, fConfidence, err := oSelf.AbstractClassifier.Best(aImage, aPokerFaceNames)
	if err != nil {
		return nil, err
	}

	return &domain.PokerFace{Name: sName, Confidence: fConfidence}, nil
}
