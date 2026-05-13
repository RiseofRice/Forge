package detection

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	namedEntityRe  = regexp.MustCompile(`&[a-zA-Z][a-zA-Z0-9]{1,10};`)
	numericEntRe   = regexp.MustCompile(`&#[0-9]{1,6};`)
	hexEntRe       = regexp.MustCompile(`&#x[0-9a-fA-F]{1,6};`)
)

// HTMLEntDetector detects HTML entity encoding.
type HTMLEntDetector struct{}

func (d *HTMLEntDetector) Name() string { return "html" }

func (d *HTMLEntDetector) Detect(data []byte) Result {
	if len(data) == 0 {
		return Result{Name: d.Name()}
	}

	s := string(data)

	named := len(namedEntityRe.FindAllString(s, -1))
	numeric := len(numericEntRe.FindAllString(s, -1))
	hexEnts := len(hexEntRe.FindAllString(s, -1))
	total := named + numeric + hexEnts

	if total == 0 {
		return Result{Name: d.Name(), Confidence: 0}
	}

	// Ratio of entity chars to total
	ratio := float64(total*4) / float64(len(s))
	if ratio > 1.0 {
		ratio = 1.0
	}

	confidence := 0.4 + ratio*0.5
	if confidence > 1.0 {
		confidence = 1.0
	}

	var kinds []string
	if named > 0 {
		kinds = append(kinds, fmt.Sprintf("%d named", named))
	}
	if numeric > 0 {
		kinds = append(kinds, fmt.Sprintf("%d numeric", numeric))
	}
	if hexEnts > 0 {
		kinds = append(kinds, fmt.Sprintf("%d hex", hexEnts))
	}

	return Result{
		Name:       d.Name(),
		Confidence: confidence,
		Details:    fmt.Sprintf("%.0f%%: HTML entities detected (%s)", confidence*100, strings.Join(kinds, ", ")),
	}
}
