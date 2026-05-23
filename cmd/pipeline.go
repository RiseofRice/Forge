package cmd

import (
	"strings"

	"github.com/RiseofRice/Forge/internal/detection"
	"github.com/RiseofRice/Forge/internal/transform"
)

// buildRegistry returns a detection Registry with built-in detectors plus
// any detectors registered by loaded plugins.
func buildRegistry() *detection.Registry {
	reg := detection.DefaultRegistry()
	for name, fn := range pluginManager.Detectors() {
		reg.Register(&detection.FuncDetector{
			DetectorName: name,
			DetectFunc:   fn,
		})
	}
	return reg
}

// pluginDecode tries plugin-registered decoders first, then falls back to
// the built-in transform.Decode.
func pluginDecode(encoding string, data []byte) ([]byte, error) {
	key := strings.ToLower(strings.TrimSpace(encoding))
	decoders := pluginManager.Decoders()
	if fn, ok := decoders[key]; ok {
		return fn(data)
	}
	return transform.Decode(encoding, data)
}

// pluginEncode tries plugin-registered encoders first, then falls back to
// the built-in transform.Encode.
func pluginEncode(encoding string, data []byte) ([]byte, error) {
	key := strings.ToLower(strings.TrimSpace(encoding))
	encoders := pluginManager.Encoders()
	if fn, ok := encoders[key]; ok {
		return fn(data)
	}
	return transform.Encode(encoding, data)
}
