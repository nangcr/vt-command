package cmdparser

import (
	"bytes"
	"github.com/chzyer/readline"
	"io"
	"io/ioutil"
	"log"
)

var _ CMDParser = &readlineImpl{}

type readlineImpl struct {
	cmd string

	logger *log.Logger
}

func NewReadline() CMDParser {
	rl := &readlineImpl{logger: log.New(io.Discard, "", log.LstdFlags)}
	return rl
}

func NewReadlineWithLogger(logger *log.Logger) CMDParser {
	rl := &readlineImpl{logger: logger}
	return rl
}

func (rl *readlineImpl) Write(r io.Reader) error {
	buffer, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	l, err := readline.NewEx(&readline.Config{Stdin: io.NopCloser(bytes.NewReader(buffer)), Stdout: ioutil.Discard})
	if err != nil {
		return err
	}

	defer l.Close()

	line, err := l.Readline()
	if err != nil {
		return err
	}

	rl.cmd = line
	return nil
}

func (rl *readlineImpl) Flush() string {
	result := rl.cmd
	rl.cmd = ""
	return result
}

func (rl *readlineImpl) Parse(r io.Reader) string {
	err := rl.Write(r)
	if err != nil {
		return ""
	}
	return rl.Flush()
}
