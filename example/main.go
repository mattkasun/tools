//nolint:mnd
package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/mattkasun/tools"
	"github.com/mattkasun/tools/logging"
)

func main() {
	overflow := math.MaxInt
	overflow += overflow / 2
	logger := logging.TextLogger(logging.TruncateSource(), logging.TimeFormat(time.DateTime))
	logger.Info("checking ...", "function", "PrettyByteSize")
	fmt.Println(999, tools.PrettyByteSize(999))
	fmt.Println(2048, tools.PrettyByteSize(2048))
	fmt.Println(1058575, tools.PrettyByteSize(1058575))
	fmt.Println(overflow, tools.PrettyByteSize(overflow))
	// Output.
	// 999 999 B
	// 2048 2.00 KiB
	// 1058575 1.01 MiB
	logger.Debug("this won't get output")
	log.Println("using std log")
	logger.Info("done")
}
