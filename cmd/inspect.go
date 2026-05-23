package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/RiseofRice/Forge/internal/analysis"
	"github.com/RiseofRice/Forge/internal/detection"
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
	Source         string             `json:"source"`
	Size           int                `json:"size_bytes"`
	MagicBytes     string             `json:"magic_bytes"`
	Entropy        float64            `json:"entropy"`
	Interpretation string             `json:"entropy_interpretation"`
	Detections     []detection.Result `json:"detections"`
}

func runInspect(cmd *cobra.Command, args []string) error {
	reg := buildRegistry()

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
			fmt.Println(separator(50))
			fmt.Printf(" %s  %s\n", bold("Source     :"), cyan(report.Source))
			fmt.Printf(" %s  %s\n", bold("Size       :"), fmt.Sprintf("%d bytes", report.Size))
			fmt.Printf(" %s  %s\n", bold("Magic bytes:"), yellow(report.MagicBytes))

			entropyStr := fmt.Sprintf("%.4f / 8.0", report.Entropy)
			entropyColored := entropyStr
			switch {
			case report.Entropy >= 7.0:
				entropyColored = red(entropyStr)
			case report.Entropy >= 5.0:
				entropyColored = yellow(entropyStr)
			default:
				entropyColored = green(entropyStr)
			}
			fmt.Printf(" %s  %s  %s\n", bold("Entropy    :"), entropyColored, dim("("+report.Interpretation+")"))
			fmt.Println(separator(50))

			if len(report.Detections) == 0 {
				fmt.Printf(" %s\n", dim("No encodings detected."))
			} else {
				fmt.Printf(" %s\n", bold("Detections :"))
				for _, r := range report.Detections {
					fmt.Printf("   %-16s %s  %s\n",
						bold(r.Name),
						confidenceBar(r.Confidence),
						r.Details,
					)
				}
			}
			fmt.Println(separator(50))
			fmt.Println()
		}
	}
	return nil
}
