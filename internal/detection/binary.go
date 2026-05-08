package detection

import "fmt"

// BinaryDetector detects common binary file formats by magic bytes.
type BinaryDetector struct{}

func (d *BinaryDetector) Name() string { return "binary" }

type binarySignature struct {
	name   string
	magic  []byte
	offset int
}

var signatures = []binarySignature{
	{name: "ELF", magic: []byte{0x7f, 0x45, 0x4c, 0x46}},
	{name: "PE/MZ", magic: []byte{0x4d, 0x5a}},
	{name: "PDF", magic: []byte{0x25, 0x50, 0x44, 0x46}}, // %PDF
	{name: "PNG", magic: []byte{0x89, 0x50, 0x4e, 0x47}}, // .PNG
	{name: "ZIP", magic: []byte{0x50, 0x4b, 0x03, 0x04}}, // PK..
	{name: "ZIP(empty)", magic: []byte{0x50, 0x4b, 0x05, 0x06}},
	{name: "JPEG", magic: []byte{0xff, 0xd8, 0xff}},
	{name: "GIF87a", magic: []byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61}},
	{name: "GIF89a", magic: []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}},
}

func (d *BinaryDetector) Detect(data []byte) Result {
	for _, sig := range signatures {
		if len(data) < sig.offset+len(sig.magic) {
			continue
		}
		match := true
		for i, b := range sig.magic {
			if data[sig.offset+i] != b {
				match = false
				break
			}
		}
		if match {
			return Result{
				Name:       d.Name(),
				Confidence: 1.0,
				Details:    fmt.Sprintf("100%%: %s magic bytes detected", sig.name),
			}
		}
	}
	return Result{Name: d.Name(), Confidence: 0}
}
