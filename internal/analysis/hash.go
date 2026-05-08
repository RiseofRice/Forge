package analysis

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

// HashResult holds a computed hash value.
type HashResult struct {
	Algorithm string `json:"algorithm"`
	Hex       string `json:"hex"`
}

// ComputeHash computes a single hash for the given algorithm.
// Supported algorithms: md5, sha1, sha256, sha512.
func ComputeHash(data []byte, algo string) (HashResult, error) {
	switch algo {
	case "md5":
		h := md5.Sum(data)
		return HashResult{Algorithm: "md5", Hex: hex.EncodeToString(h[:])}, nil
	case "sha1":
		h := sha1.Sum(data)
		return HashResult{Algorithm: "sha1", Hex: hex.EncodeToString(h[:])}, nil
	case "sha256":
		h := sha256.Sum256(data)
		return HashResult{Algorithm: "sha256", Hex: hex.EncodeToString(h[:])}, nil
	case "sha512":
		h := sha512.Sum512(data)
		return HashResult{Algorithm: "sha512", Hex: hex.EncodeToString(h[:])}, nil
	default:
		return HashResult{}, fmt.Errorf("unsupported algorithm: %s (use md5, sha1, sha256, sha512)", algo)
	}
}

// ComputeAllHashes computes hashes for all supported algorithms.
func ComputeAllHashes(data []byte) []HashResult {
	algos := []string{"md5", "sha1", "sha256", "sha512"}
	results := make([]HashResult, 0, len(algos))
	for _, algo := range algos {
		r, err := ComputeHash(data, algo)
		if err == nil {
			results = append(results, r)
		}
	}
	return results
}
