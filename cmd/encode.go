package cmd

import (
	"fmt"
	"os"

	"github.com/forgecli/forgecli/internal/transform"
	"github.com/spf13/cobra"
)

var encodeCmd = &cobra.Command{
	Use:   "encode <encoding> [file...]",
	Short: "Encode data with the specified encoding",
	Long: `Encode data using the specified encoding. Supported encodings:
  base64, base64url, hex, url, gzip, zlib`,
	Args: cobra.MinimumNArgs(1),
	RunE: runEncode,
}

func init() {
	rootCmd.AddCommand(encodeCmd)
}

func runEncode(cmd *cobra.Command, args []string) error {
	encoding := args[0]
	fileArgs := args[1:]

	inputs, err := readInputs(fileArgs)
	if err != nil {
		return err
	}

	for name, data := range inputs {
		encoded, err := transform.Encode(encoding, data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "encode error (%s) for %s: %v\n", encoding, name, err)
			os.Exit(1)
		}
		os.Stdout.Write(encoded)
	}
	return nil
}
