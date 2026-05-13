package transform

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/url"
	"strings"
)

// Decode dispatches to the appropriate decoder based on the encoding name.
func Decode(encoding string, data []byte) ([]byte, error) {
	switch strings.ToLower(strings.TrimSpace(encoding)) {
	case "base64", "b64":
		return DecodeBase64(data)
	case "base32", "b32":
		return DecodeBase32(data)
	case "hex":
		return DecodeHex(data)
	case "url", "urlencode", "urlencoding", "percent":
		return DecodeURL(data)
	case "gzip", "gz":
		return DecodeGzip(data)
	case "zlib":
		return DecodeZlib(data)
	case "jwt":
		return DecodeJWT(data)
	case "rot13":
		return DecodeROT13(data), nil
	case "html":
		return DecodeHTML(data), nil
	default:
		return nil, fmt.Errorf("unknown encoding: %s", encoding)
	}
}

// DecodeBase64 decodes standard or URL-safe base64, with or without padding.
func DecodeBase64(data []byte) ([]byte, error) {
	s := strings.TrimSpace(string(data))

	if out, err := base64.StdEncoding.DecodeString(s); err == nil {
		return out, nil
	}
	if out, err := base64.RawStdEncoding.DecodeString(s); err == nil {
		return out, nil
	}
	if out, err := base64.URLEncoding.DecodeString(s); err == nil {
		return out, nil
	}
	if out, err := base64.RawURLEncoding.DecodeString(s); err == nil {
		return out, nil
	}
	return nil, fmt.Errorf("invalid base64 data")
}

// DecodeBase32 decodes standard base32 (RFC 4648), adding padding if missing.
func DecodeBase32(data []byte) ([]byte, error) {
	s := strings.TrimSpace(strings.ToUpper(string(data)))
	for len(s)%8 != 0 {
		s += "="
	}
	out, err := base32.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("invalid base32 data: %w", err)
	}
	return out, nil
}

// DecodeHex decodes hex-encoded data, stripping an optional 0x prefix.
func DecodeHex(data []byte) ([]byte, error) {
	s := strings.TrimSpace(string(data))
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	return hex.DecodeString(s)
}

// DecodeURL decodes URL percent-encoded data.
func DecodeURL(data []byte) ([]byte, error) {
	decoded, err := url.QueryUnescape(string(data))
	if err != nil {
		return nil, fmt.Errorf("url decode: %w", err)
	}
	return []byte(decoded), nil
}

// DecodeGzip decompresses gzip data.
func DecodeGzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("gzip open: %w", err)
	}
	defer r.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("gzip read: %w", err)
	}
	return out, nil
}

// DecodeZlib decompresses zlib data.
func DecodeZlib(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("zlib open: %w", err)
	}
	defer r.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("zlib read: %w", err)
	}
	return out, nil
}

// DecodeJWT decodes a JWT and returns the payload as pretty-printed JSON.
func DecodeJWT(data []byte) ([]byte, error) {
	s := strings.TrimSpace(string(data))
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT: expected 3 parts, got %d", len(parts))
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("jwt payload decode: %w", err)
	}

	var payload interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, fmt.Errorf("jwt payload JSON: %w", err)
	}

	pretty, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("jwt payload marshal: %w", err)
	}
	return pretty, nil
}

// DecodeROT13 applies ROT13 to the input (ROT13 is its own inverse).
func DecodeROT13(data []byte) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		switch {
		case b >= 'a' && b <= 'z':
			out[i] = 'a' + (b-'a'+13)%26
		case b >= 'A' && b <= 'Z':
			out[i] = 'A' + (b-'A'+13)%26
		default:
			out[i] = b
		}
	}
	return out
}

// DecodeHTML unescapes HTML entities (&amp; &lt; &#123; etc.).
func DecodeHTML(data []byte) []byte {
	return []byte(html.UnescapeString(string(data)))
}
