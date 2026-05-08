package analysis

import "math"

// Shannon computes the Shannon entropy of data, returning a value in [0.0, 8.0].
func Shannon(data []byte) float64 {
	if len(data) == 0 {
		return 0.0
	}

	freq := make([]int, 256)
	for _, b := range data {
		freq[b]++
	}

	n := float64(len(data))
	entropy := 0.0
	for _, count := range freq {
		if count == 0 {
			continue
		}
		p := float64(count) / n
		entropy -= p * math.Log2(p)
	}
	return entropy
}

// BlockEntropy computes Shannon entropy for each block of blockSize bytes.
func BlockEntropy(data []byte, blockSize int) []float64 {
	if blockSize <= 0 || len(data) == 0 {
		return nil
	}

	var results []float64
	for i := 0; i < len(data); i += blockSize {
		end := i + blockSize
		if end > len(data) {
			end = len(data)
		}
		results = append(results, Shannon(data[i:end]))
	}
	return results
}

// InterpretEntropy returns a human-readable interpretation of an entropy value.
func InterpretEntropy(e float64) string {
	switch {
	case e < 1.0:
		return "very low (highly repetitive)"
	case e < 3.0:
		return "low (structured text or data)"
	case e < 5.0:
		return "medium (natural language or structured binary)"
	case e < 7.0:
		return "high (compressed or rich data)"
	default:
		return "very high (encrypted or random)"
	}
}
