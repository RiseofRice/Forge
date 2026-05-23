package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/RiseofRice/Forge/internal/detection"
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
	detectCmd.Flags().BoolVar(&detectAll, "all", false, "Show all detectors, including zero-confidence results")
	rootCmd.AddCommand(detectCmd)
}

func runDetect(cmd *cobra.Command, args []string) error {
	reg := buildRegistry()

	inputs, err := readInputs(args)
	if err != nil {
		return err
	}

	for name, data := range inputs {
		log.Debug().Str("source", name).Int("bytes", len(data)).Msg("detecting")

		var results []detection.Result
		if detectAll {
			results = reg.DetectAllFull(data)
		} else {
			results = reg.DetectAllParallel(data)
		}

		if outputFmt == "json" {
			type jsonResult struct {
				Source  string             `json:"source"`
				Results []detection.Result `json:"results"`
			}
			out := jsonResult{Source: name, Results: results}
			b, err := json.MarshalIndent(out, "", "  ")
			if err != nil {
				return fmt.Errorf("marshaling JSON: %w", err)
			}
			fmt.Println(string(b))
		} else {
			matches := 0
			for _, r := range results {
				if r.Confidence > 0 {
					matches++
				}
			}
			fmt.Printf("%s  %s\n", header("Source:"), cyan(name))
			fmt.Printf("%s  %d bytes", header("Size:"), len(data))
			if matches > 0 {
				fmt.Printf("  %s\n", green(fmt.Sprintf("(%d match(es))", matches)))
			} else {
				fmt.Printf("  %s\n", dim("(no matches)"))
			}
			fmt.Println(separator(50))

			if len(results) == 0 {
				fmt.Printf("  %s\n", dim("No encodings detected."))
			}
			for _, r := range results {
				details := r.Details
				if details == "" {
					details = dim("—")
				}
				fmt.Printf("  %-16s %s  %s\n",
					bold(r.Name),
					confidenceBar(r.Confidence),
					details,
				)
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
