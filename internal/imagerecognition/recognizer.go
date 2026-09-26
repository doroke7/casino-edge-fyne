// Package imagerecognition holds the image recognition logic and its gRPC service.
package imagerecognition

import (
	"bytes"
	"context"
	"image"
	_ "image/jpeg" // register decoders for image.DecodeConfig
	_ "image/png"
)

// Label is one thing found in an image.
type Label struct {
	Name       string
	Confidence float32
}

// Result is what a Recognizer reports for one image.
type Result struct {
	Labels []Label
	Width  int
	Height int
	Format string
}

// Recognizer turns encoded image bytes into a Result. Plug a real model in here.
type Recognizer interface {
	Recognize(oCtx context.Context, aImage []byte) (Result, error)
}

// BasicRecognizer only decodes the image header; it reports size and format and finds no labels.
type BasicRecognizer struct{}

func (BasicRecognizer) Recognize(_ context.Context, aImage []byte) (Result, error) {
	oConfig, sFormat, err := image.DecodeConfig(bytes.NewReader(aImage))
	if err != nil {
		return Result{}, err
	}
	return Result{Width: oConfig.Width, Height: oConfig.Height, Format: sFormat}, nil
}
