package cmd

import (
	"os"
	"strings"

	"github.com/mattn/go-isatty"
)

var colorEnabled bool

func init() {
	colorEnabled = isatty.IsTerminal(os.Stdout.Fd()) &&
		os.Getenv("NO_COLOR") == "" &&
		os.Getenv("TERM") != "dumb"
}

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiDim    = "\033[2m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiRed    = "\033[31m"
	ansiCyan   = "\033[36m"
)

func colorize(code, s string) string {
	if !colorEnabled {
		return s
	}
	return code + s + ansiReset
}

func bold(s string) string   { return colorize(ansiBold, s) }
func dim(s string) string    { return colorize(ansiDim, s) }
func green(s string) string  { return colorize(ansiGreen, s) }
func yellow(s string) string { return colorize(ansiYellow, s) }
func red(s string) string    { return colorize(ansiRed, s) }
func cyan(s string) string   { return colorize(ansiCyan, s) }

// confidenceBar returns a 20-char bar showing confidence visually.
// Colored green >= 0.75, yellow >= 0.4, red below that.
func confidenceBar(c float64) string {
	const width = 20
	filled := int(c * float64(width))
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	if !colorEnabled {
		return "[" + bar + "]"
	}
	var code string
	switch {
	case c >= 0.75:
		code = ansiGreen
	case c >= 0.40:
		code = ansiYellow
	default:
		code = ansiRed
	}
	return "[" + code + bar + ansiReset + "]"
}

// header prints a bold section header.
func header(label string) string {
	return bold(label)
}

// separator returns a dim horizontal line of the given width.
func separator(width int) string {
	return dim(strings.Repeat("─", width))
}
