package outputPortAnyModel

import (
	domain "landan-desktop-fyne/internal/domain"
)

type DieModel interface {
	// Recognize finds the dice in an encoded (JPEG or PNG) image.
	Recognize(aImage []byte) ([]*domain.Die, error)
}
