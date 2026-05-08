package detection

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Base64Detector detects standard and URL-safe base64 encoding.
type Base64Detector struct{}

func (d *Base64Detector) Name() string { return "base64" }

func (d *Base64Detector) Detect(data []byte) Result {
	if len(data) == 0 {
		return Result{Name: d.Name()}
	}

	s := strings.TrimSpace(string(data))
	if len(s) == 0 {
		return Result{Name: d.Name()}
	}

	// Try URL-safe first, then standard
	urlSafe := isBase64URLSafe(s)
	std := isBase64Standard(s)

	if !urlSafe && !std {
		return Result{Name: d.Name(), Confidence: 0}
	}

	var confidence float64
	var details string

	// Valid charset contributes 0.5
	confidence += 0.5

	// Proper length/padding contributes 0.3
	if len(s)%4 == 0 {
		confidence += 0.3
		details = "valid charset and padding"
	} else {
		details = "valid charset, no padding"
	}

	// Actually decodable contributes 0.2
	var decodeErr error
	if urlSafe {
		_, decodeErr = base64.RawURLEncoding.DecodeString(strings.TrimRight(s, "="))
		if decodeErr != nil {
			_, decodeErr = base64.URLEncoding.DecodeString(s)
		}
		details += " (url-safe)"
	} else {
		_, decodeErr = base64.StdEncoding.DecodeString(s)
		if decodeErr != nil {
			_, decodeErr = base64.RawStdEncoding.DecodeString(s)
		}
	}

	if decodeErr == nil {
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

func isBase64Standard(s string) bool {
	s = strings.TrimRight(s, "=")
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '+' || c == '/') {
			return false
		}
	}
	return len(s) > 0
}

func isBase64URLSafe(s string) bool {
	s = strings.TrimRight(s, "=")
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	return len(s) > 0
}
