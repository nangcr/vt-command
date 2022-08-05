package main

import (
	"fmt"
	"github.com/nangcr/vt-command/cmdparser"
	"os"
)

func main() {
	d := cmdparser.NewDispatcher()

	err := d.Write(os.Stdin)
	if err != nil {
		return
	}

	fmt.Print(d.Flush())
}
