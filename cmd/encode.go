package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var encodeCmd = &cobra.Command{
	Use:   "encode <encoding> [file...]",
	Short: "Encode data with the specified encoding",
	Long: `Encode data using the specified encoding. Supported encodings:
  base64, base64url, base32, hex, url, gzip, zlib, rot13, html

Plugin-registered encoders are also available.`,
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
		encoded, err := pluginEncode(encoding, data)
		if err != nil {
			return fmt.Errorf("encode error (%s) for %s: %w", encoding, name, err)
		}
		os.Stdout.Write(encoded)
	}
	return nil
}
