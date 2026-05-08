package detection

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// GzipDetector detects gzip-compressed data.
type GzipDetector struct{}

func (d *GzipDetector) Name() string { return "gzip" }

func (d *GzipDetector) Detect(data []byte) Result {
	if len(data) < 2 {
		return Result{Name: d.Name()}
	}

	// Gzip magic bytes: 0x1f 0x8b
	if data[0] != 0x1f || data[1] != 0x8b {
		return Result{Name: d.Name(), Confidence: 0}
	}

	// Magic bytes present: 0.7
	confidence := 0.7
	details := "gzip magic bytes (1f 8b)"

	// Try to decompress
	r, err := gzip.NewReader(bytes.NewReader(data))
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
