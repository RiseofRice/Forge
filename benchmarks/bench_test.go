package benchmarks

import (
	"crypto/rand"
	"testing"

	"github.com/forgecli/forgecli/internal/analysis"
	"github.com/forgecli/forgecli/internal/detection"
)

func randomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

// BenchmarkDetectAll benchmarks the full detector registry on 1KB of random data.
func BenchmarkDetectAll(b *testing.B) {
	reg := detection.DefaultRegistry()
	data := randomBytes(1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reg.DetectAll(data)
	}
}

// BenchmarkDetectAllParallel benchmarks the parallel detector registry on 1KB of random data.
func BenchmarkDetectAllParallel(b *testing.B) {
	reg := detection.DefaultRegistry()
	data := randomBytes(1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reg.DetectAllParallel(data)
	}
}

// BenchmarkBase64Detect benchmarks just the base64 detector.
func BenchmarkBase64Detect(b *testing.B) {
	d := &detection.Base64Detector{}
	data := randomBytes(1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Detect(data)
	}
}

// BenchmarkEntropy benchmarks the Shannon entropy function on 1KB of random data.
func BenchmarkEntropy(b *testing.B) {
	data := randomBytes(1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analysis.Shannon(data)
	}
}

// BenchmarkEntropyLarge benchmarks Shannon entropy on 1MB of data.
func BenchmarkEntropyLarge(b *testing.B) {
	data := randomBytes(1024 * 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analysis.Shannon(data)
	}
}
