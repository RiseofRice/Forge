# ForgeCLI

A terminal-native data transformation and analysis toolkit, inspired by [CyberChef](https://gchq.github.io/CyberChef/). Feed it raw data — files, pipes, blobs — and it detects, decodes, encodes, and analyzes without leaving your terminal.

---

## Install

```bash
make build
```

Produces a single binary: `./forge`. Move it anywhere on your `$PATH`.

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
```

Reports detected encoding(s) with a confidence score. Supports: base64, base64url, hex, URL encoding, gzip, zlib, JWT, JSON, XOR, UTF variants, and binary formats (ELF, PE, ZIP, PNG, PDF).

---

### `auto` — recursive decode

```bash
forge auto file.bin
echo "SGVsbG8=" | forge auto
```

Repeatedly detects and decodes the input until it reaches a plain-text or unrecognized form. Useful for peeling multi-layer encoded payloads.

---

### `decode` / `encode` — explicit transforms

```bash
forge decode base64 file.b64
forge decode gzip file.gz
forge decode hex file.hex
forge decode jwt token.txt

forge encode base64 file.bin
forge encode hex file.bin
forge encode url params.txt
echo "hello world" | forge encode base64
```

Supported encodings:

| Encoding   | Decode | Encode |
|------------|--------|--------|
| base64     | yes    | yes    |
| base64url  | yes    | yes    |
| hex        | yes    | yes    |
| url        | yes    | yes    |
| gzip       | yes    | yes    |
| zlib       | yes    | yes    |
| jwt        | yes    | —      |

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

Plugins implement the `Plugin` interface (`Name`, `Version`, `Register`) from `pkg/plugin/`. See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for the plugin API.

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

| Flag | Description |
|------|-------------|
| `-o, --output` | Output format: `text` (default) or `json` |
| `-v, --verbose` | Enable verbose logging |
| `--version` | Print version |

---

## How It Works

```
stdin / file
     │
     ▼
  cmd layer          (cmd/ — Cobra commands)
     │
     ▼
  detection          (internal/detection/ — 10 detectors, confidence scores)
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
