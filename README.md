# ForgeCLI

[![CI](https://github.com/RiseofRice/Forge/actions/workflows/ci.yml/badge.svg)](https://github.com/RiseofRice/Forge/actions/workflows/ci.yml)
[![Release](https://github.com/RiseofRice/Forge/actions/workflows/release.yml/badge.svg)](https://github.com/RiseofRice/Forge/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/RiseofRice/Forge)](https://goreportcard.com/report/github.com/RiseofRice/Forge)

A terminal-native data transformation and analysis toolkit, inspired by [CyberChef](https://gchq.github.io/CyberChef/). Feed it raw data — files, pipes, blobs — and it detects, decodes, encodes, and analyzes without leaving your terminal.

---

## Why Forge?

You're in a CTF, incident response, or a reverse engineering session. You have a blob of mystery data and no browser handy.

| Task | Without Forge | With Forge |
|------|--------------|------------|
| Detect encoding | `file`, guess, try `base64 -d` by hand | `forge detect payload.bin` |
| Peel multi-layer encoding | Manual trial-and-error in CyberChef | `forge auto payload.bin` |
| Decode XOR without key | Write a bruteforcer | `forge decode xor payload.bin` |
| Entropy + hashes + magic in one shot | Three separate commands | `forge inspect payload.bin` |
| Script it | Browser only | `forge detect --output json \| jq` |

**forge auto** is the killer feature: it recursively peels every encoding layer — base64 inside gzip inside hex — until it hits plaintext or an unrecognized format. `xxd`, `openssl enc`, and `jq` each handle one step at a time. Forge does the whole chain automatically.

### CTF Examples

```bash
# Challenge: you have a .bin file, no idea what it is
forge inspect challenge.bin
# → entropy 7.98 (likely encrypted/compressed), SHA-256: ..., magic: gzip

# Peel a triple-encoded payload (base64 > gzip > base64)
echo "H4sIAAAAA..." | forge auto
# → [raw] H4sI...
#   └─ [base64] .H..
#      └─ [gzip] Hello CTF{flag_here}

# XOR brute-force: recover key automatically from frequency analysis
forge decode xor encoded_payload.bin

# Hash a suspicious binary for VirusTotal lookup
forge hash --all malware.bin --output json | jq '.hashes.sha256'

# Script a full triage pipeline
forge inspect --output json suspicious.bin | jq '{entropy, magic, detections}'
```

---

## Install

### Option 1 — Download binary (no Go required)

Download the latest precompiled binary for your platform from the [Releases page](https://github.com/RiseofRice/Forge/releases):

```bash
# Linux (amd64)
curl -LO https://github.com/RiseofRice/Forge/releases/latest/download/forge_linux_amd64.tar.gz
tar -xzf forge_linux_amd64.tar.gz
sudo mv forge /usr/local/bin/

# macOS (arm64 / Apple Silicon)
curl -LO https://github.com/RiseofRice/Forge/releases/latest/download/forge_darwin_arm64.tar.gz
tar -xzf forge_darwin_arm64.tar.gz
sudo mv forge /usr/local/bin/
```

### Option 2 — go install (requires Go 1.24+)

```bash
go install github.com/RiseofRice/Forge@latest
```

### Option 3 — Build from source

```bash
git clone https://github.com/RiseofRice/Forge.git
cd Forge
make build
# binary: ./forge
```

---

## Quick Start

```bash
# What encoding is this?
forge detect suspicious.bin

# Auto-detect and peel back every layer of encoding
forge auto suspicious.bin

# Full overview: magic bytes, entropy, hashes, encoding hints
forge inspect suspicious.bin
```

---

## Commands

### `detect` — identify encoding or format

```bash
forge detect file.bin
cat file.bin | forge detect
forge detect --output json file.bin
forge detect --all file.bin        # show all detectors, including low-confidence
```

Runs all 13 detectors in parallel and reports results with confidence scores.

Supported formats: `base64`, `base32`, `hex`, `url`, `gzip`, `zlib`, `jwt`, `json`, `xor`, `rot13`, `html entities`, `utf variants`, `binary` (ELF, PE, ZIP, PNG, PDF, …).

---

### `auto` — recursive decode

```bash
forge auto file.bin
echo "SGVsbG8=" | forge auto
forge auto --max-depth 10 file.bin
```

Repeatedly detects and decodes the input until it reaches a plain-text or unrecognized form. Useful for peeling multi-layer encoded payloads. XOR-encoded data is automatically decoded using frequency analysis to find the key.

Default max depth: 5.

---

### `decode` / `encode` — explicit transforms

```bash
forge decode base64 file.b64
forge decode gzip file.gz
forge decode hex file.hex
forge decode jwt token.txt
forge decode xor payload.bin

forge encode base64 file.bin
forge encode hex file.bin
forge encode url params.txt
echo "hello world" | forge encode base64
```

Supported encodings:

| Encoding    | Decode | Encode |
|-------------|--------|--------|
| base64      | yes    | yes    |
| base64url   | yes    | yes    |
| base32      | yes    | yes    |
| hex         | yes    | yes    |
| url         | yes    | yes    |
| gzip        | yes    | yes    |
| zlib        | yes    | yes    |
| jwt         | yes    | —      |
| rot13       | yes    | yes    |
| html        | yes    | yes    |
| xor         | yes    | —      |

XOR decoding uses single-byte frequency analysis to automatically recover the key.

---

### `hash` — compute checksums

```bash
forge hash file.bin              # SHA-256 by default
forge hash --algo md5 file.bin
forge hash --all file.bin        # MD5, SHA-1, SHA-256, SHA-512
forge hash --output json file.bin
```

---

### `entropy` — Shannon entropy analysis

```bash
forge entropy file.bin
forge entropy --block-size 512 file.bin   # per-block breakdown
forge entropy --chart file.bin            # ASCII chart of entropy curve
```

High entropy (near 8.0) suggests encryption or compression. Low entropy indicates repetitive or plain-text data.

---

### `inspect` — all-in-one overview

```bash
forge inspect file.bin
```

Combines magic byte detection, entropy score, hash(es), and encoding hints into a single report. Good first command when you don't know what you're looking at.

---

### `recipe` — chained pipelines

Automate multi-step workflows with a YAML recipe file:

```yaml
name: "decode-payload"
description: "Base64 decode then decompress"
steps:
  - op: decode
    args: base64
  - op: decode
    args: gzip
```

```bash
forge recipe run decode-payload input.bin
forge recipe myrecipe.yaml input.bin
forge recipe list
```

---

### `plugins` — extend forge

```bash
forge plugins list
forge plugins info myplugin
```

Plugins implement the `Plugin` interface (`Name`, `Version`, `Register`) from `internal/plugin/`. See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for the plugin API.

---

## JSON Output

Every command supports `--output json` for scripting or piping into `jq`:

```bash
forge detect --output json file.bin | jq '.detections[0].encoding'
forge hash --output json --all file.bin | jq '.hashes'
forge entropy --output json file.bin | jq '.entropy'
```

---

## Global Flags

| Flag              | Description                              |
|-------------------|------------------------------------------|
| `-o, --output`    | Output format: `text` (default) or `json` |
| `-v, --verbose`   | Enable verbose logging                   |
| `--version`       | Print version                            |

---

## How It Works

```
stdin / file
     │
     ▼
  cmd layer          (cmd/ — Cobra commands)
     │
     ▼
  detection          (internal/detection/ — 13 detectors, parallel, confidence scores)
  transform          (internal/transform/ — decode + encode dispatchers)
  analysis           (internal/analysis/ — entropy + hashing)
     │
     ▼
  stdout (text or JSON)
```

ForgeCLI only reads and analyzes data — it never executes payloads. All commands are safe to run on untrusted input.

---

## Architecture

Full package and data-flow documentation: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT
