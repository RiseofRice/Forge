package detection

import (
	"fmt"
	"math"
)

// XORDetector detects single-byte XOR encoding by frequency analysis.
type XORDetector struct{}

func (d *XORDetector) Name() string { return "xor" }

// expectedFreq is a rough frequency table for printable ASCII / English text.
var expectedFreq = map[byte]float64{
	' ': 0.130, 'e': 0.102, 't': 0.075, 'a': 0.071, 'o': 0.068,
	'i': 0.063, 'n': 0.061, 's': 0.058, 'h': 0.058, 'r': 0.049,
	'd': 0.043, 'l': 0.035, 'u': 0.027, 'm': 0.024, 'w': 0.023,
	'c': 0.022, 'f': 0.020, 'g': 0.019, 'y': 0.018, 'p': 0.018,
}

func (d *XORDetector) Detect(data []byte) Result {
	if len(data) < 8 {
		return Result{Name: d.Name(), Confidence: 0}
	}

	bestKey, bestScore := findBestXORKey(data)
	if bestScore < 0.1 {
		return Result{Name: d.Name(), Confidence: 0}
	}

	// Normalise score to [0, 1]
	confidence := math.Min(bestScore, 1.0)

	return Result{
		Name:       d.Name(),
		Confidence: confidence,
		Details:    fmt.Sprintf("%.0f%%: best XOR key=0x%02x", confidence*100, bestKey),
	}
}

func findBestXORKey(data []byte) (byte, float64) {
	var bestKey byte
	bestScore := 0.0

	for key := 0; key < 256; key++ {
		score := scoreXOR(data, byte(key))
		if score > bestScore {
			bestScore = score
			bestKey = byte(key)
		}
	}
	return bestKey, bestScore
}

func scoreXOR(data []byte, key byte) float64 {
	freq := make(map[byte]int)
	for _, b := range data {
		freq[b^key]++
	}

	score := 0.0
	n := float64(len(data))
	for ch, weight := range expectedFreq {
		count := float64(freq[ch])
		score += (count / n) * weight
	}
	return score
}
