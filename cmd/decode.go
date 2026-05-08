package cmd

import (
	"fmt"
	"os"

	"github.com/forgecli/forgecli/internal/transform"
	"github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
	Use:   "decode <encoding> [file...]",
	Short: "Decode data with the specified encoding",
	Long: `Decode data using the specified encoding. Supported encodings:
  base64, hex, url, gzip, zlib, jwt`,
	Args: cobra.MinimumNArgs(1),
	RunE: runDecode,
}

func init() {
	rootCmd.AddCommand(decodeCmd)
}

func runDecode(cmd *cobra.Command, args []string) error {
	encoding := args[0]
	fileArgs := args[1:]

	inputs, err := readInputs(fileArgs)
	if err != nil {
		return err
	}

	for name, data := range inputs {
		decoded, err := transform.Decode(encoding, data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "decode error (%s) for %s: %v\n", encoding, name, err)
			os.Exit(1)
		}
		os.Stdout.Write(decoded)
	}
	return nil
}
