package detection

import (
	"fmt"
	"unicode/utf8"
)

// UTFDetector detects UTF encoding variants.
type UTFDetector struct{}

func (d *UTFDetector) Name() string { return "utf" }

func (d *UTFDetector) Detect(data []byte) Result {
	if len(data) == 0 {
		return Result{Name: d.Name()}
	}

	// Check BOM markers
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return Result{
			Name:       d.Name(),
			Confidence: 1.0,
			Details:    "100%: UTF-8 BOM detected",
		}
	}
	if len(data) >= 4 && data[0] == 0xFF && data[1] == 0xFE && data[2] == 0x00 && data[3] == 0x00 {
		return Result{
			Name:       d.Name(),
			Confidence: 1.0,
			Details:    "100%: UTF-32 LE BOM detected",
		}
	}
	if len(data) >= 4 && data[0] == 0x00 && data[1] == 0x00 && data[2] == 0xFE && data[3] == 0xFF {
		return Result{
			Name:       d.Name(),
			Confidence: 1.0,
			Details:    "100%: UTF-32 BE BOM detected",
		}
	}
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xFE {
		return Result{
			Name:       d.Name(),
			Confidence: 1.0,
			Details:    "100%: UTF-16 LE BOM detected",
		}
	}
	if len(data) >= 2 && data[0] == 0xFE && data[1] == 0xFF {
		return Result{
			Name:       d.Name(),
			Confidence: 1.0,
			Details:    "100%: UTF-16 BE BOM detected",
		}
	}

	// Validate UTF-8 sequences
	if utf8.Valid(data) {
		// Check ratio of multi-byte sequences to estimate likelihood
		runes := 0
		multiByteRunes := 0
		for i := 0; i < len(data); {
			r, size := utf8.DecodeRune(data[i:])
			if r == utf8.RuneError && size == 1 {
				break
			}
			runes++
			if size > 1 {
				multiByteRunes++
			}
			i += size
		}

		if multiByteRunes > 0 {
			confidence := 0.6 + float64(multiByteRunes)/float64(runes)*0.4
			return Result{
				Name:       d.Name(),
				Confidence: confidence,
				Details:    fmt.Sprintf("%.0f%%: valid UTF-8, %d multi-byte runes", confidence*100, multiByteRunes),
			}
		}
		// Pure ASCII is valid UTF-8 but not "interesting" for detection
		return Result{Name: d.Name(), Confidence: 0}
	}

	return Result{Name: d.Name(), Confidence: 0}
}
