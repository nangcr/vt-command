# vt-command

[![Go Report Card](https://goreportcard.com/badge/github.com/nangcr/vt-command?style=flat-square)](https://goreportcard.com/report/github.com/nangcr/vt-command)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/nangcr/vt-command/master/LICENSE)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/nangcr/vt-command)

A GO library for parse terminal command with ANSI escape code sequence to human-readable text.

The pkg `dispatcher` implements a dispatcher that can dispatch byte of input and finally parse to human-readable text for audit.; more information can be found here: http://www.vt100.net/emu/dec_ansi_parser.

NOTE: This library is still in development.Only support few ANSI escape code sequence now.

## Install

### Use as a CLI tool:
```bash
go install github.com/nangcr/vt-command
```
Then
```bash
echo -e '<YOUR STRING>' | vt-command > output.txt
```

### Use for development

```bash
go get -u github.com/nangcr/vt-command
```

And then import the package in your code:

```go
import "github.com/nangcr/vt-command/dispatcher"
```

### Example

An example described below is one of the use cases.

```go
package main

import (
	"fmt"
	"github.com/nangcr/vt-command/dispatcher"
	"strings"
)

func main() {
	d := dispatcher.NewDispatcher()

	err := d.Write(strings.NewReader("mysql> show databases;\u001B[9Gelect 1;\u001B[K"))
	if err != nil {
		return
	}

	fmt.Print(d.Flush())
}
```


```
mysql> select 1;
```


## License

Copyright (c) 2022-present [Nangcr](https://github.com/nangcr)

Licensed under [MIT License](./LICENSE)
