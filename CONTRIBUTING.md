# Contributing to ForgeCLI

Thank you for your interest in contributing!

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-fork/forgecli`
3. Build: `make build`
4. Run tests: `make test`

## Adding a Detector

1. Create `internal/detection/myformat.go`
2. Implement the `Detector` interface:
   ```go
   type MyFormatDetector struct{}
   func (d *MyFormatDetector) Name() string { return "myformat" }
   func (d *MyFormatDetector) Detect(data []byte) detection.Result { ... }
   ```
3. Register in `DefaultRegistry()` in `internal/detection/detector.go`
4. Add tests in `tests/detection_test.go`

## Adding an Encoder/Decoder

1. Add `EncodeXxx` / `DecodeXxx` functions to `internal/transform/encoder.go` / `decoder.go`
2. Add the case to the `Encode` / `Decode` dispatch functions
3. Add roundtrip tests

## Adding a Plugin

Implement the `plugin.Plugin` interface in `pkg/plugin/plugin.go` and call `manager.Load(myPlugin)`.

## Code Style

- Run `gofmt -w ./...` before committing
- Keep functions small and focused
- No global mutable state
- Errors go to stderr; data goes to stdout

## Pull Requests

- One feature/fix per PR
- Include tests
- Update README if adding new commands or encodings
