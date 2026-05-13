package detection

import (
	"encoding/base32"
	"fmt"
	"strings"
)

// Base32Detector detects standard base32 encoding (RFC 4648).
type Base32Detector struct{}

func (d *Base32Detector) Name() string { return "base32" }

func (d *Base32Detector) Detect(data []byte) Result {
	if len(data) == 0 {
		return Result{Name: d.Name()}
	}

	s := strings.TrimSpace(string(data))
	if len(s) == 0 {
		return Result{Name: d.Name()}
	}

	// Base32 uses A-Z and 2-7, padded with = to multiple of 8
	stripped := strings.TrimRight(s, "=")
	for _, c := range stripped {
		if !((c >= 'A' && c <= 'Z') || (c >= '2' && c <= '7')) {
			return Result{Name: d.Name(), Confidence: 0}
		}
	}

	if len(stripped) == 0 {
		return Result{Name: d.Name(), Confidence: 0}
	}

	// Valid charset: 0.5
	confidence := 0.5
	details := "valid base32 charset (A-Z, 2-7)"

	// Length multiple of 8 (with padding): 0.3
	if len(s)%8 == 0 {
		confidence += 0.3
		details += ", valid padding"
	}

	// Actually decodable: 0.2
	_, err := base32.StdEncoding.DecodeString(s)
	if err != nil {
		// Try without padding
		padded := s
		for len(padded)%8 != 0 {
			padded += "="
		}
		_, err = base32.StdEncoding.DecodeString(padded)
	}
	if err == nil {
		confidence += 0.2
		details += ", decodable"
	}

	if confidence > 1.0 {
		confidence = 1.0
	}

	return Result{
		Name:       d.Name(),
		Confidence: confidence,
		Details:    fmt.Sprintf("%.0f%%: %s", confidence*100, details),
	}
}
