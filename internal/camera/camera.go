// Package camera captures the default webcam through ffmpeg (avfoundation),
// previewing it frame by frame while recording it to an mp4 file
// and saving a still snapshot at a fixed interval.
package camera

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// expandHome turns a leading "~" into the user's home directory.
func expandHome(sPath string) (string, error) {

	if sPath != "~" && !strings.HasPrefix(sPath, "~/") {
		return sPath, nil
	}

	sHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(sHome, strings.TrimPrefix(sPath, "~")), nil
}

// NewRecordingPath returns a new timestamped mp4 path inside sDirectory, creating the folder if needed.
func NewRecordingPath(sDirectory string) (string, error) {

	sDir, err := expandHome(sDirectory)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(sDir, 0o755); err != nil {
		return "", err
	}

	return filepath.Join(sDir, "landan-"+time.Now().Format("20060102-150405")+".mp4"), nil
}

// NewSnapshotDir returns sDirectory for periodic snapshots, creating it if needed.
func NewSnapshotDir(sDirectory string) (string, error) {

	sDir, err := expandHome(sDirectory)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(sDir, 0o755); err != nil {
		return "", err
	}

	return sDir, nil
}

// Options controls what Stream writes to disk.
type Options struct {
	Width     int // capture size
	Height    int
	Framerate int    // capture frames per second
	Bitrate   string // recording bitrate, e.g. "4M"

	PreviewWidth     int // live preview size, kept small to save CPU
	PreviewHeight    int
	PreviewFramerate int

	RecordPath    string        // mp4 file to record into
	SnapshotDir   string        // folder for snapshots; empty disables snapshots
	SnapshotEvery time.Duration // interval between snapshots; the first one is taken right away
}

// readJPEG returns the next JPEG from a concatenated MJPEG stream.
// Inside JPEG entropy data every 0xFF is stuffed as FF 00, so the bytes FF D9 only ever mean "end of image".
func readJPEG(oReader *bufio.Reader) ([]byte, error) {

	var aData []byte
	for {
		aChunk, err := oReader.ReadBytes(0xD9)
		aData = append(aData, aChunk...)
		if err != nil {
			return nil, err
		}
		if n := len(aData); n >= 2 && aData[n-2] == 0xFF {
			return aData, nil
		}
	}
}

// Stream reads the default webcam with ffmpeg (avfoundation).
// One ffmpeg process feeds up to three outputs: an MJPEG stream on stdout for the live preview,
// an H.264 mp4 file, and periodic full-size JPEG snapshots.
// Cancelling oCtx asks ffmpeg to quit so the mp4 is finalized.
func Stream(oCtx context.Context, oOptions Options, fnFrame func(image.Image)) error {

	aArgs := []string{
		"-hide_banner", "-loglevel", "error",
		"-f", "avfoundation",
		"-framerate", strconv.Itoa(oOptions.Framerate),
		"-video_size", fmt.Sprintf("%dx%d", oOptions.Width, oOptions.Height),
		"-pixel_format", "nv12", // what the hardware encoder takes directly, avoids a conversion
		"-i", "default:none",
		// output 1: live preview, kept small and slow to save CPU
		"-vf", fmt.Sprintf("scale=%d:%d", oOptions.PreviewWidth, oOptions.PreviewHeight),
		"-r", strconv.Itoa(oOptions.PreviewFramerate),
		"-f", "image2pipe", "-vcodec", "mjpeg", "-q:v", "5", "-",
		// output 2: recording at full size (fragmented mp4 stays playable even if the process dies)
		"-c:v", "h264_videotoolbox", "-b:v", oOptions.Bitrate,
		"-movflags", "+frag_keyframe+empty_moov", "-y", oOptions.RecordPath,
	}

	if oOptions.SnapshotDir != "" && oOptions.SnapshotEvery >= time.Second {
		// output 3: one full-size still every SnapshotEvery, named by time
		aArgs = append(aArgs,
			"-vf", fmt.Sprintf("fps=1/%d", int(oOptions.SnapshotEvery.Seconds())),
			"-q:v", "2", "-f", "image2", "-strftime", "1",
			filepath.Join(oOptions.SnapshotDir, "snapshot-%Y%m%d-%H%M%S.jpg"),
		)
	}

	oCmd := exec.CommandContext(oCtx, "ffmpeg", aArgs...)

	oStdin, err := oCmd.StdinPipe()
	if err != nil {
		return err
	}
	oStdout, err := oCmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Instead of killing ffmpeg, send "q" so it flushes and closes the mp4 properly.
	oCmd.Cancel = func() error {
		_, err := io.WriteString(oStdin, "q")
		return err
	}
	oCmd.WaitDelay = 5 * time.Second

	if err := oCmd.Start(); err != nil {
		return fmt.Errorf("無法啟動 ffmpeg(需要先安裝 ffmpeg): %w", err)
	}

	oReader := bufio.NewReaderSize(oStdout, 1<<20)
	for {
		aData, err := readJPEG(oReader)
		if err != nil {
			break
		}
		oFrame, err := jpeg.Decode(bytes.NewReader(aData))
		if err != nil {
			continue
		}
		fnFrame(oFrame)
	}

	oCmd.Wait()
	if oCtx.Err() != nil {
		return nil
	}
	return fmt.Errorf("無法讀取攝影機,請確認已允許攝影機權限,且沒有被其他程式佔用")
}
