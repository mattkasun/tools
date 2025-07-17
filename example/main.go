//nolint:mnd
package main

import (
	"fmt"

	"github.com/mattkasun/tools"
)

func main() {
	fmt.Println(999, tools.PrettyByteSize(999))
	fmt.Println(2048, tools.PrettyByteSize(2048))
	fmt.Println(1058575, tools.PrettyByteSize(1058575))
	// Output.
	// 999 999 B
	// 2048 2.00 KiB
	// 1058575 1.01 MiB
}
