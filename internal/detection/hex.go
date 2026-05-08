package detection

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// HexDetector detects hex-encoded data.
type HexDetector struct{}

func (d *HexDetector) Name() string { return "hex" }

func (d *HexDetector) Detect(data []byte) Result {
	if len(data) == 0 {
		return Result{Name: d.Name()}
	}

	s := strings.TrimSpace(string(data))
	// Remove optional 0x prefix
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")

	if len(s) == 0 {
		return Result{Name: d.Name()}
	}

	// Check charset
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return Result{Name: d.Name(), Confidence: 0}
		}
	}

	var confidence float64
	var details string

	// Valid charset: 0.5
	confidence += 0.5
	details = "valid hex charset"

	// Even length: 0.3
	if len(s)%2 == 0 {
		confidence += 0.3
		details += ", even length"
	} else {
		return Result{
			Name:       d.Name(),
			Confidence: confidence,
			Details:    fmt.Sprintf("%.0f%%: %s (odd length)", confidence*100, details),
		}
	}

	// Decodable: 0.2
	_, err := hex.DecodeString(s)
	if err == nil {
		confidence += 0.2
		details += ", decodable"
	}

	return Result{
		Name:       d.Name(),
		Confidence: confidence,
		Details:    fmt.Sprintf("%.0f%%: %s", confidence*100, details),
	}
}
