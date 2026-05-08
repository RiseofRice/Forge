package detection

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

// ZlibDetector detects zlib-compressed data.
type ZlibDetector struct{}

func (d *ZlibDetector) Name() string { return "zlib" }

func (d *ZlibDetector) Detect(data []byte) Result {
	if len(data) < 2 {
		return Result{Name: d.Name()}
	}

	// zlib magic byte combinations
	hasMagic := (data[0] == 0x78 && (data[1] == 0x9c || data[1] == 0x01 || data[1] == 0xda || data[1] == 0x5e))
	if !hasMagic {
		return Result{Name: d.Name(), Confidence: 0}
	}

	// Magic bytes present: 0.8
	confidence := 0.8
	details := fmt.Sprintf("zlib magic bytes (%02x %02x)", data[0], data[1])

	// Try to decompress
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err == nil {
		_, readErr := io.ReadAll(r)
		r.Close()
		if readErr == nil {
			confidence = 1.0
			details += ", decompressible"
		}
	}

	return Result{
		Name:       d.Name(),
		Confidence: confidence,
		Details:    fmt.Sprintf("%.0f%%: %s", confidence*100, details),
	}
}
