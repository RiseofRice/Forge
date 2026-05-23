package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/RiseofRice/Forge/internal/transform"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Recipe struct {
	Name        string       `yaml:"name" json:"name"`
	Description string       `yaml:"description" json:"description"`
	Steps       []RecipeStep `yaml:"steps" json:"steps"`
}

type RecipeStep struct {
	Op   string `yaml:"op" json:"op"`
	Args string `yaml:"args" json:"args"`
}

var recipeCmd = &cobra.Command{
	Use:   "recipe [recipe-file] [input...]",
	Short: "Run a pipeline of operations from a YAML recipe file",
	Long:  `Load and execute a YAML recipe file that defines a pipeline of encode/decode operations.`,
	RunE:  runRecipe,
}

var recipeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List built-in recipes",
	RunE:  runRecipeList,
}

var recipeRunCmd = &cobra.Command{
	Use:   "run <name>",
	Short: "Run a built-in recipe by name",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runRecipeRun,
}

var builtinRecipes = []Recipe{
	{
		Name:        "decode-b64-then-gunzip",
		Description: "Decode base64 then decompress gzip",
		Steps: []RecipeStep{
			{Op: "decode", Args: "base64"},
			{Op: "decode", Args: "gzip"},
		},
	},
	{
		Name:        "encode-gzip-then-b64",
		Description: "Compress with gzip then encode as base64",
		Steps: []RecipeStep{
			{Op: "encode", Args: "gzip"},
			{Op: "encode", Args: "base64"},
		},
	},
}

func init() {
	recipeCmd.AddCommand(recipeListCmd)
	recipeCmd.AddCommand(recipeRunCmd)
	rootCmd.AddCommand(recipeCmd)
}

func runRecipe(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	recipeFile := args[0]
	fileArgs := args[1:]

	f, err := os.Open(recipeFile)
	if err != nil {
		return fmt.Errorf("opening recipe file: %w", err)
	}
	defer f.Close()

	var recipe Recipe
	if err := yaml.NewDecoder(f).Decode(&recipe); err != nil {
		return fmt.Errorf("parsing recipe: %w", err)
	}

	inputs, err := readInputs(fileArgs)
	if err != nil {
		return err
	}

	for name, data := range inputs {
		result, err := applyRecipe(recipe, data)
		if err != nil {
			return fmt.Errorf("recipe error for %s: %w", name, err)
		}
		os.Stdout.Write(result)
	}
	return nil
}

func runRecipeList(cmd *cobra.Command, args []string) error {
	if outputFmt == "json" {
		b, err := json.MarshalIndent(builtinRecipes, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		fmt.Println(string(b))
	} else {
		fmt.Println("Built-in recipes:")
		for _, r := range builtinRecipes {
			fmt.Printf("  %-30s %s\n", r.Name, r.Description)
		}
	}
	return nil
}

func runRecipeRun(cmd *cobra.Command, args []string) error {
	name := args[0]
	fileArgs := args[1:]

	var recipe *Recipe
	for i := range builtinRecipes {
		if builtinRecipes[i].Name == name {
			recipe = &builtinRecipes[i]
			break
		}
	}
	if recipe == nil {
		return fmt.Errorf("unknown recipe: %s", name)
	}

	inputs, err := readInputs(fileArgs)
	if err != nil {
		return err
	}

	for name, data := range inputs {
		result, err := applyRecipe(*recipe, data)
		if err != nil {
			return fmt.Errorf("recipe error for %s: %w", name, err)
		}
		os.Stdout.Write(result)
	}
	return nil
}

func applyRecipe(recipe Recipe, data []byte) ([]byte, error) {
	current := data
	for i, step := range recipe.Steps {
		var err error
		switch step.Op {
		case "decode":
			current, err = transform.Decode(step.Args, current)
		case "encode":
			current, err = transform.Encode(step.Args, current)
		default:
			return nil, fmt.Errorf("step %d: unknown op %q", i+1, step.Op)
		}
		if err != nil {
			return nil, fmt.Errorf("step %d (%s %s): %w", i+1, step.Op, step.Args, err)
		}
	}
	return current, nil
}
