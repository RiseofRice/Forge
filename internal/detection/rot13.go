package detection

import (
	"fmt"
	"unicode"
)

// ROT13Detector detects ROT13-encoded text by checking whether
// applying ROT13 produces output with higher English letter frequency.
type ROT13Detector struct{}

func (d *ROT13Detector) Name() string { return "rot13" }

func (d *ROT13Detector) Detect(data []byte) Result {
	if len(data) < 8 {
		return Result{Name: d.Name(), Confidence: 0}
	}

	letterCount := 0
	for _, b := range data {
		if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') {
			letterCount++
		}
	}
	if letterCount < 4 {
		return Result{Name: d.Name(), Confidence: 0}
	}

	original := scoreEnglish(data)
	rotated := scoreEnglish(applyROT13(data))

	if rotated <= original {
		return Result{Name: d.Name(), Confidence: 0}
	}

	// How much better does ROT13 score?
	improvement := (rotated - original) / (original + 0.001)
	confidence := 0.0
	switch {
	case improvement > 0.5:
		confidence = 0.9
	case improvement > 0.2:
		confidence = 0.7
	case improvement > 0.1:
		confidence = 0.5
	default:
		confidence = 0.3
	}

	return Result{
		Name:       d.Name(),
		Confidence: confidence,
		Details:    fmt.Sprintf("%.0f%%: ROT13 improves English score by %.0f%%", confidence*100, improvement*100),
	}
}

func applyROT13(data []byte) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		switch {
		case b >= 'a' && b <= 'z':
			out[i] = 'a' + (b-'a'+13)%26
		case b >= 'A' && b <= 'Z':
			out[i] = 'A' + (b-'A'+13)%26
		default:
			out[i] = b
		}
	}
	return out
}

var englishFreq = map[rune]float64{
	'e': 0.130, 't': 0.091, 'a': 0.082, 'o': 0.075, 'i': 0.070,
	'n': 0.067, 's': 0.063, 'h': 0.061, 'r': 0.060, 'd': 0.043,
	'l': 0.040, 'u': 0.028, 'm': 0.024, 'w': 0.023, 'f': 0.022,
	'g': 0.020, 'y': 0.020, 'p': 0.019, 'b': 0.015, 'v': 0.010,
}

func scoreEnglish(data []byte) float64 {
	freq := make(map[rune]int)
	total := 0
	for _, b := range data {
		r := unicode.ToLower(rune(b))
		if r >= 'a' && r <= 'z' {
			freq[r]++
			total++
		}
	}
	if total == 0 {
		return 0
	}
	score := 0.0
	for ch, weight := range englishFreq {
		score += (float64(freq[ch]) / float64(total)) * weight
	}
	return score
}
