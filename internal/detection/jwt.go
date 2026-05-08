package detection

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// JWTDetector detects JSON Web Tokens.
type JWTDetector struct{}

func (d *JWTDetector) Name() string { return "jwt" }

func (d *JWTDetector) Detect(data []byte) Result {
	s := strings.TrimSpace(string(data))

	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return Result{Name: d.Name(), Confidence: 0}
	}

	// 3 parts separated by dots: 0.4
	confidence := 0.4
	details := "3-part dot-separated structure"

	// Decode header and check for "alg"
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Result{
			Name:       d.Name(),
			Confidence: confidence,
			Details:    fmt.Sprintf("%.0f%%: %s, invalid header encoding", confidence*100, details),
		}
	}

	var header map[string]interface{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return Result{
			Name:       d.Name(),
			Confidence: confidence,
			Details:    fmt.Sprintf("%.0f%%: %s, invalid header JSON", confidence*100, details),
		}
	}

	if _, ok := header["alg"]; ok {
		confidence += 0.4
		details += ", valid JWT header with alg"
	}

	// Decode payload and check it's valid JSON
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err == nil {
		var payload map[string]interface{}
		if json.Unmarshal(payloadBytes, &payload) == nil {
			confidence += 0.2
			details += ", valid JSON payload"
		}
	}

	return Result{
		Name:       d.Name(),
		Confidence: confidence,
		Details:    fmt.Sprintf("%.0f%%: %s", confidence*100, details),
	}
}
