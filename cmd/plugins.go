package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/forgecli/forgecli/pkg/plugin"
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
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
	} else {
		if len(plugins) == 0 {
			fmt.Println("No plugins loaded.")
			return nil
		}
		fmt.Println("Loaded plugins:")
		for _, p := range plugins {
			fmt.Printf("  %-20s v%s\n", p.Name(), p.Version())
		}
	}
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
				b, _ := json.MarshalIndent(out, "", "  ")
				fmt.Println(string(b))
			} else {
				fmt.Printf("Name   : %s\n", p.Name())
				fmt.Printf("Version: %s\n", p.Version())
			}
			return nil
		}
	}
	fmt.Fprintf(os.Stderr, "plugin not found: %s\n", name)
	os.Exit(1)
	return nil
}
