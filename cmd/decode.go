package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
	Use:   "decode <encoding> [file...]",
	Short: "Decode data with the specified encoding",
	Long: `Decode data using the specified encoding. Supported encodings:
  base64, base32, hex, url, gzip, zlib, jwt, rot13, html

Plugin-registered decoders are also available.`,
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
		decoded, err := pluginDecode(encoding, data)
		if err != nil {
			return fmt.Errorf("decode error (%s) for %s: %w", encoding, name, err)
		}
		os.Stdout.Write(decoded)
	}
	return nil
}
