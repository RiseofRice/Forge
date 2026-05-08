package cmd

import (
	"fmt"
	"strings"

	"github.com/forgecli/forgecli/internal/detection"
	"github.com/forgecli/forgecli/internal/transform"
	"github.com/spf13/cobra"
)

var autoMaxDepth int

var autoCmd = &cobra.Command{
	Use:   "auto [file...]",
	Short: "Auto-detect and decode input data recursively",
	Long:  `Automatically detect encodings and decode them recursively until raw data or max depth is reached.`,
	RunE:  runAuto,
}

func init() {
	autoCmd.Flags().IntVar(&autoMaxDepth, "max-depth", 5, "Maximum recursion depth")
	rootCmd.AddCommand(autoCmd)
}

type treeNode struct {
	encoding string
	data     []byte
	children []*treeNode
}

func runAuto(cmd *cobra.Command, args []string) error {
	reg := detection.DefaultRegistry()

	inputs, err := readInputs(args)
	if err != nil {
		return err
	}

	for name, data := range inputs {
		fmt.Printf("Source: %s\n", name)
		node := buildTree(reg, data, autoMaxDepth, 0)
		printTree(node, 0)
		fmt.Println()
	}
	return nil
}

func buildTree(reg *detection.Registry, data []byte, maxDepth, depth int) *treeNode {
	node := &treeNode{encoding: "raw", data: data}
	if depth >= maxDepth {
		return node
	}

	results := reg.DetectAllParallel(data)
	for _, r := range results {
		if r.Confidence < 0.5 {
			continue
		}
		decoded, err := transform.Decode(r.Name, data)
		if err != nil {
			continue
		}
		if string(decoded) == string(data) {
			continue
		}
		child := buildTree(reg, decoded, maxDepth, depth+1)
		child.encoding = r.Name
		node.children = append(node.children, child)
		break // take the best match only
	}
	return node
}

func printTree(node *treeNode, depth int) {
	indent := strings.Repeat("  ", depth)
	preview := previewBytes(node.data, 64)
	fmt.Printf("%s[%s] %s\n", indent, node.encoding, preview)
	for _, child := range node.children {
		printTree(child, depth+1)
	}
}

func previewBytes(data []byte, maxLen int) string {
	s := string(data)
	if len(s) > maxLen {
		s = s[:maxLen] + "..."
	}
	// Replace non-printable chars
	var b strings.Builder
	for _, r := range s {
		if r < 32 || r > 126 {
			b.WriteRune('.')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
