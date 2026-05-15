# ForgeCLI Architecture

## Overview

ForgeCLI is structured around three core layers: **detection**, **analysis**, and **transformation**. All business logic lives under `internal/` and `pkg/`, with the CLI surface in `cmd/`.

```
cmd/          — Cobra CLI commands (entry points for the user)
internal/     — Core business logic (not exported)
  detection/  — Format/encoding detectors
  analysis/   — Entropy and hash computation
  transform/  — Encoders and decoders
pkg/          — Public plugin API
  plugin/     — Plugin interface and manager
plugins/      — Drop-in plugin implementations
recipes/      — YAML recipe files
```

## Layer Descriptions

### cmd/

Each subcommand (`detect`, `encode`, `decode`, `hash`, `entropy`, `inspect`, `auto`, `recipe`, `plugins`) is a separate file. Commands read from stdin or files, delegate to the internal packages, then format output as text or JSON.

`pipeline.go` provides two shared helpers used by `auto`, `detect`, and `inspect`:
- `buildRegistry()` — constructs a `detection.Registry` from all built-in detectors plus any plugin-registered detectors
- `pluginDecode(encoding, data)` — dispatches to plugin decoders first, then falls back to `transform.Decode`

### internal/detection/

- `detector.go` — `Detector` interface, `Result` struct, `Registry` with sequential (`DetectAll`) and parallel (`DetectAllParallel`) execution, and `DefaultRegistry()` factory
- `funcdetector.go` — adapter that wraps a plain function into the `Detector` interface (used by the plugin system)

Individual detectors (each implements `Detect([]byte) Result`, returns confidence in [0.0, 1.0]):

| File          | Name       | Method                                      |
|---------------|------------|---------------------------------------------|
| `base64.go`   | base64     | character set + length heuristics           |
| `base32.go`   | base32     | character set + padding heuristics          |
| `hex.go`      | hex        | character set + even-length check           |
| `gzip.go`     | gzip       | magic bytes `1f 8b`                         |
| `zlib.go`     | zlib       | magic byte patterns                         |
| `jwt.go`      | jwt        | three-part dot structure + header decoding  |
| `json.go`     | json       | structural JSON parsing                     |
| `urlenc.go`   | urlencode  | percent-encoded sequence density            |
| `htmlent.go`  | html       | HTML entity pattern matching                |
| `xor.go`      | xor        | single-byte frequency analysis over 256 keys|
| `rot13.go`    | rot13      | English frequency score improvement         |
| `utf.go`      | utf        | BOM detection and multi-byte sequence check |
| `binary.go`   | binary     | magic bytes for ELF, PE, ZIP, PNG, PDF, …   |

### internal/analysis/

- `entropy.go` — Shannon entropy (overall and per-block) plus interpretation
- `hash.go` — MD5, SHA-1, SHA-256, SHA-512 via standard library

### internal/transform/

- `decoder.go` — `Decode(encoding, data)` dispatcher and individual `DecodeXxx` functions  
  Supported: `base64`, `base32`, `hex`, `url`, `gzip`, `zlib`, `jwt`, `rot13`, `html`, `xor`
- `encoder.go` — `Encode(encoding, data)` dispatcher and individual `EncodeXxx` functions  
  Supported: `base64`, `base64url`, `base32`, `hex`, `url`, `gzip`, `zlib`, `rot13`, `html`

XOR decoding (`DecodeXOR`) uses single-byte frequency analysis over all 256 candidate keys and applies the one that best matches expected English character distribution.

### pkg/plugin/

Public API for third-party plugins. A plugin implements `Plugin` (Name, Version, Register) and uses `Registry` to attach detectors, decoders, and encoders at startup.

## Data Flow

```
stdin / file
     │
     ▼
  cmd layer
     │  reads input, selects subcommand
     ▼
  detection.Registry.DetectAllParallel(data)
     │  runs all 13 detectors concurrently
     │  returns []Result sorted by confidence desc
     ▼
  transform.Decode / transform.Encode
     │  dispatches by encoding name
     ▼
  stdout  (text or JSON)
```

For `forge auto`, the flow is recursive: each decoded output is fed back into `DetectAllParallel` until no high-confidence match remains or `--max-depth` is reached.

## Design Principles

- **No global state** — dependencies are passed explicitly
- **Unix-native** — reads stdin, writes stdout, errors to stderr
- **JSON output** — all commands support `--output json` for pipeline use
- **Modular** — new detectors/encoders slot in without touching existing code
- **Safe** — analysis only; no execution of payloads
