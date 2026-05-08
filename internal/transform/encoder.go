package transform

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

// Encode dispatches to the appropriate encoder based on the encoding name.
func Encode(encoding string, data []byte) ([]byte, error) {
	switch strings.ToLower(strings.TrimSpace(encoding)) {
	case "base64", "b64":
		return EncodeBase64(data), nil
	case "base64url", "b64url":
		return EncodeBase64URL(data), nil
	case "hex":
		return EncodeHex(data), nil
	case "url", "urlencode", "urlencoding", "percent":
		return EncodeURL(data), nil
	case "gzip", "gz":
		return EncodeGzip(data)
	case "zlib":
		return EncodeZlib(data)
	default:
		return nil, fmt.Errorf("unknown encoding: %s", encoding)
	}
}

// EncodeBase64 encodes data as standard base64 with padding.
func EncodeBase64(data []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(data))
}

// EncodeBase64URL encodes data as URL-safe base64 without padding.
func EncodeBase64URL(data []byte) []byte {
	return []byte(base64.RawURLEncoding.EncodeToString(data))
}

// EncodeHex encodes data as lowercase hex.
func EncodeHex(data []byte) []byte {
	return []byte(hex.EncodeToString(data))
}

// EncodeURL percent-encodes data suitable for use in a query string.
func EncodeURL(data []byte) []byte {
	return []byte(url.QueryEscape(string(data)))
}

// EncodeGzip compresses data using gzip.
func EncodeGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("gzip write: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("gzip close: %w", err)
	}
	return buf.Bytes(), nil
}

// EncodeZlib compresses data using zlib.
func EncodeZlib(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("zlib write: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("zlib close: %w", err)
	}
	return buf.Bytes(), nil
}
