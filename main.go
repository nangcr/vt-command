package main

import (
	"fmt"
	"github.com/nangcr/vt-command/dispatcher"
	"os"
)

func main() {
	d := dispatcher.NewDispatcher()

	err := d.Write(os.Stdin)
	if err != nil {
		return
	}

	fmt.Print(d.Flush())
}
