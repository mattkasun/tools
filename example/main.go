//nolint:mnd,forbidigo
package main

import (
	"fmt"
	"time"

	"github.com/mattkasun/tools"
	"github.com/mattkasun/tools/logging"
)

func main() {
	log := logging.TextLogger(logging.TruncateSource(), logging.TimeFormat(time.DateTime))
	log.Info("checking ...", "function", "PrettyByteSize")
	fmt.Println(999, tools.PrettyByteSize(999))
	fmt.Println(2048, tools.PrettyByteSize(2048))
	fmt.Println(1058575, tools.PrettyByteSize(1058575))
	// Output.
	// 999 999 B
	// 2048 2.00 KiB
	// 1058575 1.01 MiB
	log.Debug("this won't get output")
	discard := logging.DiscardLogger()
	discard.Info("this won't get output")
	log.Info("done")
}
