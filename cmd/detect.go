package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/forgecli/forgecli/internal/detection"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	detectDepth int
	detectAll   bool
)

var detectCmd = &cobra.Command{
	Use:   "detect [file...]",
	Short: "Detect encoding/format of input data",
	Long:  `Detect the encoding or format of the provided input. Reads from stdin if no files are given.`,
	RunE:  runDetect,
}

func init() {
	detectCmd.Flags().IntVar(&detectDepth, "depth", 3, "Recursion depth for nested encoding detection")
	detectCmd.Flags().BoolVar(&detectAll, "all", false, "Show all detectors, not just matches")
	rootCmd.AddCommand(detectCmd)
}

func runDetect(cmd *cobra.Command, args []string) error {
	reg := detection.DefaultRegistry()

	inputs, err := readInputs(args)
	if err != nil {
		return err
	}

	for name, data := range inputs {
		log.Debug().Str("source", name).Int("bytes", len(data)).Msg("detecting")
		results := reg.DetectAllParallel(data)

		if outputFmt == "json" {
			type jsonResult struct {
				Source  string             `json:"source"`
				Results []detection.Result `json:"results"`
			}
			out := jsonResult{Source: name, Results: results}
			b, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(b))
		} else {
			fmt.Printf("Source: %s\n", name)
			if len(results) == 0 {
				fmt.Println("  No encodings detected.")
			}
			for _, r := range results {
				if !detectAll && r.Confidence == 0 {
					continue
				}
				fmt.Printf("  %-20s confidence=%.2f  %s\n", r.Name, r.Confidence, r.Details)
			}
			fmt.Println()
		}
	}
	return nil
}

// readInputs reads data from files (or stdin if args is empty).
// Returns map[source_name]data.
func readInputs(args []string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	if len(args) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("reading stdin: %w", err)
		}
		result["stdin"] = data
		return result, nil
	}
	for _, path := range args {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening %s: %v\n", path, err)
			continue
		}
		data, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
			continue
		}
		result[path] = data
	}
	return result, nil
}
