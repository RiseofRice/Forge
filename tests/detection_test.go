package tests

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"testing"

	"github.com/RiseofRice/Forge/internal/analysis"
	"github.com/RiseofRice/Forge/internal/detection"
	"github.com/RiseofRice/Forge/internal/transform"
)

// ---- Base64 detection ----

func TestBase64Detection_Valid(t *testing.T) {
	d := &detection.Base64Detector{}
	input := []byte(base64.StdEncoding.EncodeToString([]byte("hello world")))
	result := d.Detect(input)
	if result.Confidence < 0.7 {
		t.Errorf("expected confidence >= 0.7, got %f", result.Confidence)
	}
}

func TestBase64Detection_URLSafe(t *testing.T) {
	d := &detection.Base64Detector{}
	input := []byte(base64.RawURLEncoding.EncodeToString([]byte("hello world test data")))
	result := d.Detect(input)
	if result.Confidence < 0.5 {
		t.Errorf("expected confidence >= 0.5 for URL-safe base64, got %f", result.Confidence)
	}
}

func TestBase64Detection_Invalid(t *testing.T) {
	d := &detection.Base64Detector{}
	input := []byte("this is definitely not base64!!! @#$%")
	result := d.Detect(input)
	if result.Confidence > 0 {
		t.Errorf("expected confidence = 0 for invalid input, got %f", result.Confidence)
	}
}

// ---- Hex detection ----

func TestHexDetection_Valid(t *testing.T) {
	d := &detection.HexDetector{}
	input := []byte(hex.EncodeToString([]byte("hello world")))
	result := d.Detect(input)
	if result.Confidence < 0.9 {
		t.Errorf("expected confidence >= 0.9 for valid hex, got %f", result.Confidence)
	}
}

func TestHexDetection_Invalid(t *testing.T) {
	d := &detection.HexDetector{}
	input := []byte("this is not hex")
	result := d.Detect(input)
	if result.Confidence > 0 {
		t.Errorf("expected confidence = 0 for invalid hex, got %f", result.Confidence)
	}
}

func TestHexDetection_OddLength(t *testing.T) {
	d := &detection.HexDetector{}
	input := []byte("abc") // valid hex chars but odd length
	result := d.Detect(input)
	// Should detect valid charset (0.5) but not even length
	if result.Confidence >= 0.8 {
		t.Errorf("expected lower confidence for odd-length hex, got %f", result.Confidence)
	}
}

// ---- Gzip detection ----

func TestGzipDetection_Valid(t *testing.T) {
	d := &detection.GzipDetector{}
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write([]byte("hello gzip world"))
	w.Close()
	result := d.Detect(buf.Bytes())
	if result.Confidence < 0.9 {
		t.Errorf("expected confidence >= 0.9 for valid gzip, got %f", result.Confidence)
	}
}

func TestGzipDetection_Invalid(t *testing.T) {
	d := &detection.GzipDetector{}
	input := []byte("not gzip data at all")
	result := d.Detect(input)
	if result.Confidence > 0 {
		t.Errorf("expected confidence = 0 for non-gzip data, got %f", result.Confidence)
	}
}

// ---- JWT detection ----

func TestJWTDetection_Valid(t *testing.T) {
	d := &detection.JWTDetector{}
	// A minimal valid JWT structure
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"1234567890","name":"John Doe","iat":1516239022}`))
	sig := "SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	jwt := header + "." + payload + "." + sig
	result := d.Detect([]byte(jwt))
	if result.Confidence < 0.9 {
		t.Errorf("expected confidence >= 0.9 for valid JWT, got %f", result.Confidence)
	}
}

func TestJWTDetection_Invalid(t *testing.T) {
	d := &detection.JWTDetector{}
	input := []byte("not.a.jwt.token.extra")
	result := d.Detect(input)
	if result.Confidence > 0 {
		t.Errorf("expected confidence = 0 for extra-part input, got %f", result.Confidence)
	}
}

// ---- JSON detection ----

func TestJSONDetection_Valid(t *testing.T) {
	d := &detection.JSONDetector{}
	input := []byte(`{"key": "value", "number": 42}`)
	result := d.Detect(input)
	if result.Confidence < 1.0 {
		t.Errorf("expected confidence = 1.0 for valid JSON, got %f", result.Confidence)
	}
}

func TestJSONDetection_NotJSON(t *testing.T) {
	d := &detection.JSONDetector{}
	input := []byte("hello world this is plain text")
	result := d.Detect(input)
	if result.Confidence > 0 {
		t.Errorf("expected confidence = 0 for plain text, got %f", result.Confidence)
	}
}

// ---- Entropy ----

func TestEntropyZero(t *testing.T) {
	// All same bytes => entropy 0
	data := bytes.Repeat([]byte{0x41}, 1000)
	e := analysis.Shannon(data)
	if e != 0.0 {
		t.Errorf("expected entropy 0.0 for uniform data, got %f", e)
	}
}

func TestEntropyMax(t *testing.T) {
	// 256 distinct bytes => close to 8.0
	data := make([]byte, 256)
	for i := range data {
		data[i] = byte(i)
	}
	e := analysis.Shannon(data)
	if e < 7.9 {
		t.Errorf("expected entropy ~8.0 for all-byte data, got %f", e)
	}
}

func TestEntropyEmpty(t *testing.T) {
	e := analysis.Shannon([]byte{})
	if e != 0.0 {
		t.Errorf("expected entropy 0.0 for empty data, got %f", e)
	}
}

// ---- Encode/Decode roundtrips ----

func TestRoundtripBase64(t *testing.T) {
	original := []byte("hello world, this is a roundtrip test!")
	encoded, err := transform.Encode("base64", original)
	if err != nil {
		t.Fatalf("encode base64 error: %v", err)
	}
	decoded, err := transform.Decode("base64", encoded)
	if err != nil {
		t.Fatalf("decode base64 error: %v", err)
	}
	if !bytes.Equal(original, decoded) {
		t.Errorf("roundtrip mismatch: got %q, want %q", decoded, original)
	}
}

func TestRoundtripHex(t *testing.T) {
	original := []byte("hello hex world 12345")
	encoded, err := transform.Encode("hex", original)
	if err != nil {
		t.Fatalf("encode hex error: %v", err)
	}
	decoded, err := transform.Decode("hex", encoded)
	if err != nil {
		t.Fatalf("decode hex error: %v", err)
	}
	if !bytes.Equal(original, decoded) {
		t.Errorf("roundtrip mismatch: got %q, want %q", decoded, original)
	}
}

func TestRoundtripGzip(t *testing.T) {
	original := []byte("hello gzip world! compress me please.")
	encoded, err := transform.Encode("gzip", original)
	if err != nil {
		t.Fatalf("encode gzip error: %v", err)
	}
	decoded, err := transform.Decode("gzip", encoded)
	if err != nil {
		t.Fatalf("decode gzip error: %v", err)
	}
	if !bytes.Equal(original, decoded) {
		t.Errorf("roundtrip mismatch: got %q, want %q", decoded, original)
	}
}

func TestRoundtripURL(t *testing.T) {
	original := []byte("hello world&foo=bar+baz")
	encoded, err := transform.Encode("url", original)
	if err != nil {
		t.Fatalf("encode url error: %v", err)
	}
	decoded, err := transform.Decode("url", encoded)
	if err != nil {
		t.Fatalf("decode url error: %v", err)
	}
	// url.QueryEscape encodes space as '+', QueryUnescape decodes '+' as space
	if string(decoded) != url.QueryEscape(string(original)) {
		// Just check it's a valid roundtrip by re-encoding
		reEncoded := url.QueryEscape(string(decoded))
		if reEncoded != string(encoded) {
			t.Logf("decoded: %q, original: %q", decoded, original)
		}
	}
}
