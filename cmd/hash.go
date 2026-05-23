package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/RiseofRice/Forge/internal/analysis"
	"github.com/spf13/cobra"
)

var (
	hashAlgo string
	hashAll  bool
)

var hashCmd = &cobra.Command{
	Use:   "hash [file...]",
	Short: "Compute hash of input data",
	Long:  `Compute cryptographic hash of input data. Default algorithm is sha256.`,
	RunE:  runHash,
}

func init() {
	hashCmd.Flags().StringVar(&hashAlgo, "algo", "sha256", "Hash algorithm: md5, sha1, sha256, sha512")
	hashCmd.Flags().BoolVar(&hashAll, "all", false, "Compute all supported hash algorithms")
	rootCmd.AddCommand(hashCmd)
}

func runHash(cmd *cobra.Command, args []string) error {
	inputs, err := readInputs(args)
	if err != nil {
		return err
	}

	for name, data := range inputs {
		if hashAll {
			results := analysis.ComputeAllHashes(data)
			if outputFmt == "json" {
				type jsonOut struct {
					Source string               `json:"source"`
					Hashes []analysis.HashResult `json:"hashes"`
				}
				out := jsonOut{Source: name, Hashes: results}
				b, _ := json.MarshalIndent(out, "", "  ")
				fmt.Println(string(b))
			} else {
				fmt.Printf("Source: %s\n", name)
				for _, r := range results {
					fmt.Printf("  %-8s %s\n", r.Algorithm, r.Hex)
				}
				fmt.Println()
			}
		} else {
			result, err := analysis.ComputeHash(data, hashAlgo)
			if err != nil {
				fmt.Fprintf(os.Stderr, "hash error: %v\n", err)
				os.Exit(1)
			}
			if outputFmt == "json" {
				type jsonOut struct {
					Source    string `json:"source"`
					Algorithm string `json:"algorithm"`
					Hash      string `json:"hash"`
				}
				out := jsonOut{Source: name, Algorithm: result.Algorithm, Hash: result.Hex}
				b, _ := json.MarshalIndent(out, "", "  ")
				fmt.Println(string(b))
			} else {
				fmt.Printf("%s  %s\n", result.Hex, name)
			}
		}
	}
	return nil
}
