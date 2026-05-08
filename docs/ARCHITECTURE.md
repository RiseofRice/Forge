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
plugins/      — Drop-in plugin implementations (future)
recipes/      — YAML recipe files
```

## Layer Descriptions

### cmd/
Each subcommand (`detect`, `encode`, `decode`, `hash`, `entropy`, `inspect`, `auto`, `recipe`, `plugins`) is a separate file. Commands read from stdin or files, delegate to the internal packages, then format output as text or JSON.

### internal/detection/
- `detector.go` — `Detector` interface, `Result` struct, `Registry` (sequential and parallel execution), `DefaultRegistry()` factory
- Individual detector files (`base64.go`, `hex.go`, `gzip.go`, `zlib.go`, `jwt.go`, `json.go`, `urlenc.go`, `xor.go`, `utf.go`, `binary.go`) each implement `Detector` with a `Detect([]byte) Result` method that returns a confidence score in [0.0, 1.0]

### internal/analysis/
- `entropy.go` — Shannon entropy (overall and per-block) plus interpretation
- `hash.go` — MD5, SHA-1, SHA-256, SHA-512 via standard library

### internal/transform/
- `decoder.go` — `Decode(encoding, data)` dispatcher and individual `DecodeXxx` functions
- `encoder.go` — `Encode(encoding, data)` dispatcher and individual `EncodeXxx` functions

### pkg/plugin/
Public API for third-party plugins. A plugin implements `Plugin` (Name, Version, Register) and uses `Registry` to attach detectors, decoders, and encoders.

## Data Flow

```
stdin / file  →  cmd layer  →  detection / transform / analysis  →  stdout (text or JSON)
```

## Design Principles

- **No global state** — dependencies are passed explicitly
- **Unix-native** — reads stdin, writes stdout, errors to stderr
- **JSON output** — all commands support `--output json` for pipeline use
- **Modular** — new detectors/encoders slot in without touching existing code
- **Safe** — analysis only; no execution of payloads
