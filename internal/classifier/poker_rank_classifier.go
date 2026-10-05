package classifier

import (
	"fmt"

	bootstrap "landan-desktop-fyne/bootstrap"
	domain "landan-desktop-fyne/internal/domain"
	pkgInference "landan-desktop-fyne/pkg/inference"
)

// 撲克牌的點數，順序見 pk-studio-ir-model 的 cfg/classify/poker/rank-data.yaml。
var aPokerRankNames = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

// PokerRankClassifier 用 onnx 跑 撲克牌的點數，onnx 的細節都在 AbstractClassifier。
type PokerRankClassifier struct {
	*AbstractClassifier
}

// NewPokerRankClassifier 從 config/inference.yaml 的 classify.poker.rank 讀模型路徑。
func NewPokerRankClassifier(oInference *pkgInference.Inference) (*PokerRankClassifier, error) {
	if bootstrap.CONFIG.INFERENCE.CLASSIFY.POKER.RANK.ONNX == "" {
		return nil, fmt.Errorf("inference.classify.poker.rank.onnx is empty (is config/inference.yaml filled in? run from the project root)")
	}

	oAbstractClassifier, err := NewAbstractClassifier(oInference, bootstrap.CONFIG.INFERENCE.CLASSIFY.POKER.RANK.ONNX, bootstrap.CONFIG.INFERENCE.CLASSIFY.POKER.RANK.OPENVINO)
	if err != nil {
		return nil, err
	}

	return &PokerRankClassifier{AbstractClassifier: oAbstractClassifier}, nil
}

func (oSelf *PokerRankClassifier) Classify(aImage []byte) (*domain.PokerRank, error) {
	sName, fConfidence, err := oSelf.AbstractClassifier.Best(aImage, aPokerRankNames)
	if err != nil {
		return nil, err
	}

	return &domain.PokerRank{Name: sName, Confidence: fConfidence}, nil
}
