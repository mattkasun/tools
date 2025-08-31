// Package tools implements utility routines
package tools

import (
	"fmt"
	"math"
	"os"

	"golang.org/x/term"
)

const (
	// Escape is ansi terminal escape sequence.
	Escape = "\x1b["
	// Reset is ansi termnial code for reset.
	Reset = Escape + "0m"
	// Green is ansi terminal code for green text.
	Green = Escape + "32m"
)

// PrettyByteSize formats a byte size (int) into a human-readable colored string using binary prefixes (KiB, MiB, etc.).
func PrettyByteSize(b int) string {
	var response string
	float := float64(b)
	for i, unit := range []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"} {
		if math.Abs(float) < 1024.0 { //nolint:mnd
			if UseColour() {
				unit = Green + unit + Reset
			}
			if i == 0 {
				response = fmt.Sprintf("%1.0f %s", float, unit)
				break
			}
			response = fmt.Sprintf("%3.2f %s", float, unit)
			break
		}
		float /= 1024.0
	}
	return response
}

// UseColour checks if output is a terminal or whether colour output has been suppresed in env.
func UseColour() bool {
	_, ok := os.LookupEnv("NO_COLOUR")
	if term.IsTerminal(int(os.Stdout.Fd())) && !ok {
		return true
	}
	return false
}
