package cmdparser

import (
	"io"
)

// CMDParser is the interface that wraps the methods of the dispatcherImpl
type CMDParser interface {
	Parse(r io.Reader) string
	Write(r io.Reader) error
	Flush() string
}
