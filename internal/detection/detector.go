package detection

import (
	"sort"
	"sync"
)

// Result holds the output of a single detector.
type Result struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"` // 0.0 - 1.0
	Details    string  `json:"details"`
}

// Detector is the interface implemented by all format detectors.
type Detector interface {
	Name() string
	Detect(data []byte) Result
}

// Registry holds all registered detectors.
type Registry struct {
	mu        sync.RWMutex
	detectors []Detector
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register adds a detector to the registry.
func (r *Registry) Register(d Detector) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.detectors = append(r.detectors, d)
}

// DetectAll runs all detectors sequentially and returns results with Confidence > 0,
// sorted by confidence descending.
func (r *Registry) DetectAll(data []byte) []Result {
	r.mu.RLock()
	detectors := make([]Detector, len(r.detectors))
	copy(detectors, r.detectors)
	r.mu.RUnlock()

	results := make([]Result, 0, len(detectors))
	for _, d := range detectors {
		res := d.Detect(data)
		if res.Confidence > 0 {
			results = append(results, res)
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Confidence > results[j].Confidence
	})
	return results
}

// DetectAllParallel runs all detectors in parallel using goroutines,
// returning only results with Confidence > 0, sorted descending.
func (r *Registry) DetectAllParallel(data []byte) []Result {
	r.mu.RLock()
	detectors := make([]Detector, len(r.detectors))
	copy(detectors, r.detectors)
	r.mu.RUnlock()

	type indexed struct {
		idx int
		res Result
	}
	ch := make(chan indexed, len(detectors))

	var wg sync.WaitGroup
	for i, d := range detectors {
		wg.Add(1)
		go func(idx int, det Detector) {
			defer wg.Done()
			ch <- indexed{idx: idx, res: det.Detect(data)}
		}(i, d)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	results := make([]Result, 0, len(detectors))
	for item := range ch {
		if item.res.Confidence > 0 {
			results = append(results, item.res)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Confidence > results[j].Confidence
	})
	return results
}

// DetectAllFull runs all detectors in parallel and returns every result
// (including zero-confidence), sorted by confidence descending.
// Used by the --all flag to show all detectors regardless of match.
func (r *Registry) DetectAllFull(data []byte) []Result {
	r.mu.RLock()
	detectors := make([]Detector, len(r.detectors))
	copy(detectors, r.detectors)
	r.mu.RUnlock()

	type indexed struct {
		idx int
		res Result
	}
	ch := make(chan indexed, len(detectors))

	var wg sync.WaitGroup
	for i, d := range detectors {
		wg.Add(1)
		go func(idx int, det Detector) {
			defer wg.Done()
			ch <- indexed{idx: idx, res: det.Detect(data)}
		}(i, d)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	results := make([]Result, 0, len(detectors))
	for item := range ch {
		results = append(results, item.res)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Confidence > results[j].Confidence
	})
	return results
}

// DefaultRegistry returns a Registry pre-loaded with all built-in detectors.
func DefaultRegistry() *Registry {
	reg := NewRegistry()
	reg.Register(&Base64Detector{})
	reg.Register(&Base32Detector{})
	reg.Register(&HexDetector{})
	reg.Register(&GzipDetector{})
	reg.Register(&ZlibDetector{})
	reg.Register(&JWTDetector{})
	reg.Register(&JSONDetector{})
	reg.Register(&URLEncDetector{})
	reg.Register(&HTMLEntDetector{})
	reg.Register(&XORDetector{})
	reg.Register(&ROT13Detector{})
	reg.Register(&UTFDetector{})
	reg.Register(&BinaryDetector{})
	return reg
}
