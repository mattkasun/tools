//nolint:testpackage
package tools

import (
	"math"
	"testing"
)

func Test_prettyByteSize(t *testing.T) {
	overflow := math.MaxInt
	overflow++
	t.Parallel()
	tests := []struct {
		name string
		args int
		want string
	}{
		{"Zero bytes", 0, "0 B"},
		{"Bytes under 1 KiB", 512, "512 B"},
		{"Exactly 1 KiB", 1024, "1.00 KiB"},
		{"1.5 KiB", 1536, "1.50 KiB"},
		{"Exactly 1 MiB", 1024 * 1024, "1.00 MiB"},
		{"Exactly 1 GiB", 1024 * 1024 * 1024, "1.00 GiB"},
		{"Large value (1 TiB)", 1 << 40, "1.00 TiB"},
		{"Large value (1 PiB)", 1 << 50, "1.00 PiB"},
		{"Large value (1 EiB)", 1 << 60, "1.00 EiB"},
		{"Very large (beyond EiB)", math.MaxInt, "8.00 EiB"}, // Max handled before overflow fallback
		{"Overflow", overflow, "-8.00 EiB"},
		{"Negative size", -1536, "-1.50 KiB"}, // still shows absolute value
	}
	for _, tt := range tests {
		// tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := PrettyByteSize(tt.args); got != tt.want {
				t.Errorf("prettyByteSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
