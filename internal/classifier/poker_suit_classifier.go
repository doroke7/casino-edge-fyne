package classifier

import (
	"fmt"

	bootstrap "landan-desktop-fyne/bootstrap"
	domain "landan-desktop-fyne/internal/domain"
	pkgInference "landan-desktop-fyne/pkg/inference"
)

// 撲克牌的花色，順序見 pk-studio-ir-model 的 cfg/classify/poker/suit-data.yaml。
var aPokerSuitNames = []string{"Spade", "Heart", "Diamond", "Club"}

// PokerSuitClassifier 用 onnx 跑 撲克牌的花色，onnx 的細節都在 AbstractClassifier。
type PokerSuitClassifier struct {
	*AbstractClassifier
}

// NewPokerSuitClassifier 從 config/inference.yaml 的 classify.poker.suit 讀模型路徑。
func NewPokerSuitClassifier(oInference *pkgInference.Inference) (*PokerSuitClassifier, error) {
	if bootstrap.CONFIG.INFERENCE.CLASSIFY.POKER.SUIT.ONNX == "" {
		return nil, fmt.Errorf("inference.classify.poker.suit.onnx is empty (is config/inference.yaml filled in? run from the project root)")
	}

	oAbstractClassifier, err := NewAbstractClassifier(oInference, bootstrap.CONFIG.INFERENCE.CLASSIFY.POKER.SUIT.ONNX, bootstrap.CONFIG.INFERENCE.CLASSIFY.POKER.SUIT.OPENVINO)
	if err != nil {
		return nil, err
	}

	return &PokerSuitClassifier{AbstractClassifier: oAbstractClassifier}, nil
}

func (oSelf *PokerSuitClassifier) Classify(aImage []byte) (*domain.PokerSuit, error) {
	sName, fConfidence, err := oSelf.AbstractClassifier.Best(aImage, aPokerSuitNames)
	if err != nil {
		return nil, err
	}

	return &domain.PokerSuit{Name: sName, Confidence: fConfidence}, nil
}
