package tools_test

import (
	"math"
	"testing"

	"github.com/Kairum-Labs/should"
	"github.com/mattkasun/tools"
)

func TestPrettyByteSize(t *testing.T) {
	overflow := math.MaxInt
	overflow++
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
	t.Run("colour output", func(t *testing.T) {
		expected := "B"
		if tools.UseColour() {
			expected = tools.Green + "B" + tools.Reset
		}
		should.BeEqual(t, tools.PrettyByteSize(642), "642 "+expected)
		// t.Log(tools.PrettyByteSize(642))
	})
	for _, tt := range tests {
		// tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("NO_COLOUR", "true")
			should.BeEqual(t, tools.PrettyByteSize(tt.args), tt.want)
		})
	}
}
