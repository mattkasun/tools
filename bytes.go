// Package tools implements utility routines
package tools

import (
	"fmt"
	"math"

	"github.com/fatih/color"
)

// PrettyByteSize formats a byte size (int) into a human-readable colored string using binary prefixes (KiB, MiB, etc.).
func PrettyByteSize(b int) string {
	float := float64(b)
	for i, unit := range []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"} {
		if math.Abs(float) < 1024.0 { //nolint:mnd
			if i == 0 {
				return fmt.Sprintf("%1.0f %s", float, color.GreenString(unit))
			}
			return fmt.Sprintf("%3.2f %s", float, color.GreenString(unit))
		}
		float /= 1024.0
	}
	return color.RedString("out of range")
}
