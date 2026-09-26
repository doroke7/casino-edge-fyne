package camera

import (
	"bufio"
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"testing"
)

func encodeJPEG(t *testing.T, nWidth, nHeight int) []byte {
	t.Helper()

	oImage := image.NewRGBA(image.Rect(0, 0, nWidth, nHeight))
	for y := 0; y < nHeight; y++ {
		for x := 0; x < nWidth; x++ {
			oImage.Set(x, y, color.RGBA{R: uint8(x * 7), G: uint8(y * 11), B: 200, A: 255})
		}
	}

	var oBuffer bytes.Buffer
	if err := jpeg.Encode(&oBuffer, oImage, nil); err != nil {
		t.Fatal(err)
	}
	return oBuffer.Bytes()
}

func TestReadJPEGSplitsConcatenatedStream(t *testing.T) {
	aFirst := encodeJPEG(t, 64, 48)
	aSecond := encodeJPEG(t, 32, 24)
	aThird := encodeJPEG(t, 16, 16)

	oReader := bufio.NewReader(bytes.NewReader(bytes.Join([][]byte{aFirst, aSecond, aThird}, nil)))

	for i, aWant := range [][]byte{aFirst, aSecond, aThird} {
		aGot, err := readJPEG(oReader)
		if err != nil {
			t.Fatalf("frame %d: %v", i, err)
		}
		if !bytes.Equal(aGot, aWant) {
			t.Fatalf("frame %d: got %d bytes, want %d", i, len(aGot), len(aWant))
		}
		if _, err := jpeg.Decode(bytes.NewReader(aGot)); err != nil {
			t.Fatalf("frame %d does not decode: %v", i, err)
		}
	}

	if _, err := readJPEG(oReader); err != io.EOF {
		t.Fatalf("want io.EOF at end of stream, got %v", err)
	}
}
