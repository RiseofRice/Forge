package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/forgecli/forgecli/internal/analysis"
	"github.com/forgecli/forgecli/internal/detection"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [file...]",
	Short: "Inspect file type, magic bytes, entropy, and encoding hints",
	Long:  `Show detailed information about the input: magic bytes, detected format, entropy, size, and encoding hints.`,
	RunE:  runInspect,
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

type inspectReport struct {
	Source       string             `json:"source"`
	Size         int                `json:"size_bytes"`
	MagicBytes   string             `json:"magic_bytes"`
	Entropy      float64            `json:"entropy"`
	Interpretation string           `json:"entropy_interpretation"`
	Detections   []detection.Result `json:"detections"`
}

func runInspect(cmd *cobra.Command, args []string) error {
	reg := detection.DefaultRegistry()

	inputs, err := readInputs(args)
	if err != nil {
		return err
	}

	for name, data := range inputs {
		entropy := analysis.Shannon(data)
		interp := analysis.InterpretEntropy(entropy)
		results := reg.DetectAllParallel(data)

		magic := ""
		if len(data) >= 4 {
			magic = fmt.Sprintf("%02x %02x %02x %02x", data[0], data[1], data[2], data[3])
		} else {
			for i, b := range data {
				if i > 0 {
					magic += " "
				}
				magic += fmt.Sprintf("%02x", b)
			}
		}

		report := inspectReport{
			Source:         name,
			Size:           len(data),
			MagicBytes:     magic,
			Entropy:        entropy,
			Interpretation: interp,
			Detections:     results,
		}

		if outputFmt == "json" {
			b, _ := json.MarshalIndent(report, "", "  ")
			fmt.Println(string(b))
		} else {
			fmt.Printf("Source      : %s\n", report.Source)
			fmt.Printf("Size        : %d bytes\n", report.Size)
			fmt.Printf("Magic bytes : %s\n", report.MagicBytes)
			fmt.Printf("Entropy     : %.4f (%s)\n", report.Entropy, report.Interpretation)
			fmt.Println("Detections  :")
			for _, r := range report.Detections {
				fmt.Printf("  %-20s %.2f  %s\n", r.Name, r.Confidence, r.Details)
			}
			fmt.Println()
		}
	}
	return nil
}
