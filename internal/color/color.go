package color

import (
	"fmt"
	"math"
	"os"
	"strconv"
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

// Basic ANSI colors as RGB for nearest-color fallback
type basicColor struct {
	code string
	r, g, b uint8
}

var basicColors = []basicColor{
	{"\033[90m", 80, 80, 80},   // gray
	{"\033[91m", 255, 80, 80},  // red
	{"\033[92m", 80, 255, 80},  // green
	{"\033[93m", 255, 220, 80}, // yellow
	{"\033[34m", 60, 80, 200},  // blue
	{"\033[95m", 200, 80, 200}, // magenta/pink
	{"\033[96m", 80, 220, 220}, // cyan
	{"\033[94m", 80, 120, 255}, // bright blue
	{"\033[35m", 160, 60, 180}, // purple/violet
	{"\033[33m", 200, 130, 30}, // orange
	{"\033[32m", 40, 180, 40},  // dark green
	{"\033[36m", 0, 150, 150},  // teal
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

func truecolorSupported() bool {
	ct := os.Getenv("COLORTERM")
	return ct == "truecolor" || ct == "24bit"
}

func parseHex(hex string) (r, g, b uint8, ok bool) {
	hex = strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(hex) != 6 {
		return 0, 0, 0, false
	}
	rv, err1 := strconv.ParseUint(hex[0:2], 16, 8)
	gv, err2 := strconv.ParseUint(hex[2:4], 16, 8)
	bv, err3 := strconv.ParseUint(hex[4:6], 16, 8)
	if err1 != nil || err2 != nil || err3 != nil {
		return 0, 0, 0, false
	}
	return uint8(rv), uint8(gv), uint8(bv), true
}

func nearestAnsi(r, g, b uint8) string {
	best, bestDist := "", math.MaxFloat64
	for _, bc := range basicColors {
		dr := float64(int(r) - int(bc.r))
		dg := float64(int(g) - int(bc.g))
		db := float64(int(b) - int(bc.b))
		dist := dr*dr + dg*dg + db*db
		if dist < bestDist {
			bestDist = dist
			best = bc.code
		}
	}
	return best
}

func Apply(text, colorName string) string {
	if !Supported() {
		return text
	}

	c := strings.TrimSpace(colorName)

	// Handle hex colors
	if strings.HasPrefix(c, "#") {
		r, g, b, ok := parseHex(c)
		if !ok {
			return text
		}
		var code string
		if truecolorSupported() {
			code = fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
		} else {
			code = nearestAnsi(r, g, b)
		}
		return code + text + Reset
	}

	// Named color
	code, ok := ansiMap[strings.ToLower(c)]
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
