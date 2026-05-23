package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RiseofRice/Forge/internal/analysis"
	"github.com/spf13/cobra"
)

var (
	entropyBlockSize int
	entropyChart     bool
)

var entropyCmd = &cobra.Command{
	Use:   "entropy [file...]",
	Short: "Calculate Shannon entropy of input data",
	Long:  `Calculate Shannon entropy (0.0 - 8.0) of input data. Higher values indicate more randomness.`,
	RunE:  runEntropy,
}

func init() {
	entropyCmd.Flags().IntVar(&entropyBlockSize, "block-size", 256, "Block size for block entropy analysis")
	entropyCmd.Flags().BoolVar(&entropyChart, "chart", false, "Show ASCII chart of block entropy")
	rootCmd.AddCommand(entropyCmd)
}

func runEntropy(cmd *cobra.Command, args []string) error {
	inputs, err := readInputs(args)
	if err != nil {
		return err
	}

	for name, data := range inputs {
		overall := analysis.Shannon(data)
		interp := analysis.InterpretEntropy(overall)
		blocks := analysis.BlockEntropy(data, entropyBlockSize)

		if outputFmt == "json" {
			type jsonOut struct {
				Source         string    `json:"source"`
				Entropy        float64   `json:"entropy"`
				Interpretation string    `json:"interpretation"`
				Blocks         []float64 `json:"block_entropy,omitempty"`
			}
			out := jsonOut{
				Source:         name,
				Entropy:        overall,
				Interpretation: interp,
				Blocks:         blocks,
			}
			b, err := json.MarshalIndent(out, "", "  ")
			if err != nil {
				return fmt.Errorf("marshaling JSON: %w", err)
			}
			fmt.Println(string(b))
		} else {
			fmt.Printf("Source : %s\n", name)
			fmt.Printf("Entropy: %.4f / 8.0  [%s]\n", overall, interp)
			if entropyChart && len(blocks) > 0 {
				fmt.Println("Block entropy chart:")
				printEntropyChart(blocks)
			}
			fmt.Println()
		}
	}
	return nil
}

func printEntropyChart(blocks []float64) {
	const chartHeight = 8
	const barChar = "█"
	const emptyChar = " "

	// build columns
	cols := make([]int, len(blocks))
	for i, e := range blocks {
		cols[i] = int((e / 8.0) * float64(chartHeight))
	}

	// print top to bottom
	for row := chartHeight; row >= 1; row-- {
		fmt.Printf("%d |", row)
		for _, h := range cols {
			if h >= row {
				fmt.Print(barChar)
			} else {
				fmt.Print(emptyChar)
			}
		}
		fmt.Println()
	}
	fmt.Printf("  +%s\n", strings.Repeat("-", len(blocks)))
	fmt.Printf("   blocks (size=%d)\n", entropyBlockSize)
}
