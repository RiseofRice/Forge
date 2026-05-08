package detection

import (
	"fmt"
	"net/url"
)

// URLEncDetector detects URL-encoded (percent-encoded) data.
type URLEncDetector struct{}

func (d *URLEncDetector) Name() string { return "url" }

func (d *URLEncDetector) Detect(data []byte) Result {
	if len(data) == 0 {
		return Result{Name: d.Name()}
	}

	s := string(data)
	total := len(s)
	encodedCount := 0

	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) && isHexByte(s[i+1]) && isHexByte(s[i+2]) {
			encodedCount++
			i += 2
		}
	}

	if encodedCount == 0 {
		// Check for + (space encoding) or already-encoded strings
		_, err := url.QueryUnescape(s)
		if err == nil && s != "" {
			// No percent encoding found
			return Result{Name: d.Name(), Confidence: 0}
		}
		return Result{Name: d.Name(), Confidence: 0}
	}

	ratio := float64(encodedCount*3) / float64(total)
	if ratio > 1.0 {
		ratio = 1.0
	}

	confidence := 0.3 + ratio*0.7
	if confidence > 1.0 {
		confidence = 1.0
	}

	return Result{
		Name:       d.Name(),
		Confidence: confidence,
		Details:    fmt.Sprintf("%.0f%%: %d percent-encoded sequences found", confidence*100, encodedCount),
	}
}

func isHexByte(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
