package detection

// FuncDetector wraps a plain confidence function as a Detector.
// Used by the plugin bridge to integrate plugin-registered detectors.
type FuncDetector struct {
	DetectorName string
	DetectFunc   func([]byte) float64
}

func (f *FuncDetector) Name() string { return f.DetectorName }

func (f *FuncDetector) Detect(data []byte) Result {
	c := f.DetectFunc(data)
	if c > 1.0 {
		c = 1.0
	}
	if c < 0 {
		c = 0
	}
	return Result{Name: f.DetectorName, Confidence: c}
}
