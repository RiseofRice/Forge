package detection

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONDetector detects JSON-formatted data.
type JSONDetector struct{}

func (d *JSONDetector) Name() string { return "json" }

func (d *JSONDetector) Detect(data []byte) Result {
	if len(data) == 0 {
		return Result{Name: d.Name()}
	}

	s := strings.TrimSpace(string(data))
	if len(s) == 0 {
		return Result{Name: d.Name()}
	}

	// Try to unmarshal
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err == nil {
		return Result{
			Name:       d.Name(),
			Confidence: 1.0,
			Details:    "100%: valid JSON",
		}
	}

	// Looks like JSON (starts with { or [)
	if s[0] == '{' || s[0] == '[' {
		return Result{
			Name:       d.Name(),
			Confidence: 0.3,
			Details:    fmt.Sprintf("30%%: starts with %c but invalid JSON", s[0]),
		}
	}

	return Result{Name: d.Name(), Confidence: 0}
}
