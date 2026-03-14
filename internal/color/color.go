package color

import (
	"os"
	"strings"
)

const (
	Reset = "\033[0m"
	Bold  = "\033[1m"
	Dim   = "\033[2m"
)

var ansiMap = map[string]string{
	"cyan":          "\033[96m",
	"blue":          "\033[34m",
	"green":         "\033[32m",
	"red":           "\033[91m",
	"purple":        "\033[35m",
	"orange":        "\033[33m",
	"teal":          "\033[36m",
	"indigo":        "\033[94m",
	"pink":          "\033[95m",
	"gold":          "\033[93m",
	"amber":         "\033[93m",
	"neon-green":    "\033[92m",
	"neon-cyan":     "\033[96m",
	"metallic-blue": "\033[34m",
	"yellow":        "\033[93m",
	"violet":        "\033[35m",
	"rose":          "\033[91m",
	"lime":          "\033[92m",
	"gray":          "\033[90m",
	"fuchsia":       "\033[95m",
}

func Supported() bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func Apply(text, colorName string) string {
	if !Supported() {
		return text
	}
	code, ok := ansiMap[strings.ToLower(strings.TrimSpace(colorName))]
	if !ok {
		return text
	}
	return code + text + Reset
}

func ApplyBold(text string) string {
	if !Supported() {
		return text
	}
	return Bold + text + Reset
}

func ApplyDim(text string) string {
	if !Supported() {
		return text
	}
	return Dim + text + Reset
}
