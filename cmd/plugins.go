package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/RiseofRice/Forge/internal/plugin"
	"github.com/spf13/cobra"
)

var pluginManager = plugin.NewManager()

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Manage and list ForgeCLI plugins",
	Long:  `List and inspect loaded ForgeCLI plugins.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var pluginsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all loaded plugins",
	RunE:  runPluginsList,
}

var pluginsInfoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show details about a specific plugin",
	Args:  cobra.ExactArgs(1),
	RunE:  runPluginsInfo,
}

func init() {
	pluginsCmd.AddCommand(pluginsListCmd)
	pluginsCmd.AddCommand(pluginsInfoCmd)
	rootCmd.AddCommand(pluginsCmd)
}

func runPluginsList(cmd *cobra.Command, args []string) error {
	plugins := pluginManager.List()
	if outputFmt == "json" {
		type pInfo struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}
		out := make([]pInfo, len(plugins))
		for i, p := range plugins {
			out[i] = pInfo{Name: p.Name(), Version: p.Version()}
		}
		b, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		fmt.Println(string(b))
		return nil
	}

	if len(plugins) == 0 {
		fmt.Println(dim("No plugins loaded."))
		return nil
	}

	detectors := pluginManager.Detectors()
	decoders := pluginManager.Decoders()
	encoders := pluginManager.Encoders()

	fmt.Printf("%s\n", bold("Loaded plugins:"))
	for _, p := range plugins {
		fmt.Printf("  %-20s %s\n", bold(p.Name()), cyan("v"+p.Version()))
	}
	fmt.Println()
	fmt.Printf("  Detectors : %d registered\n", len(detectors))
	fmt.Printf("  Decoders  : %d registered\n", len(decoders))
	fmt.Printf("  Encoders  : %d registered\n", len(encoders))

	return nil
}

func runPluginsInfo(cmd *cobra.Command, args []string) error {
	name := args[0]
	plugins := pluginManager.List()
	for _, p := range plugins {
		if p.Name() == name {
			if outputFmt == "json" {
				type pInfo struct {
					Name    string `json:"name"`
					Version string `json:"version"`
				}
				out := pInfo{Name: p.Name(), Version: p.Version()}
				b, err := json.MarshalIndent(out, "", "  ")
				if err != nil {
					return fmt.Errorf("marshaling JSON: %w", err)
				}
				fmt.Println(string(b))
			} else {
				fmt.Printf("%s  %s\n", bold("Name   :"), p.Name())
				fmt.Printf("%s  %s\n", bold("Version:"), p.Version())
			}
			return nil
		}
	}
	return fmt.Errorf("plugin not found: %s", name)
}
