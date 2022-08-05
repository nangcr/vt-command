package cmdparser

import (
	"fmt"
	"github.com/danielgatis/go-vte/vtparser"
	"io"
	"log"
)

var _ CMDParser = &dispatcherImpl{}

type dispatcherImpl struct {
	cmd []rune // command buffer
	cur int    // cursor position

	parser *vtparser.Parser

	logger *log.Logger
}

func NewDispatcher() CMDParser {
	dispatcher := &dispatcherImpl{cmd: []rune{}, logger: log.New(io.Discard, "", log.LstdFlags)}
	parser := vtparser.New(
		dispatcher.print,
		dispatcher.execute,
		dispatcher.put,
		dispatcher.unhook,
		dispatcher.hook,
		dispatcher.oscDispatch,
		dispatcher.csiDispatch,
		dispatcher.escDispatch,
	)
	dispatcher.parser = parser
	return dispatcher
}

func NewDispatcherWithLogger(logger *log.Logger) *dispatcherImpl {
	dispatcher := &dispatcherImpl{cmd: []rune{}, logger: logger}
	parser := vtparser.New(
		dispatcher.print,
		dispatcher.execute,
		dispatcher.put,
		dispatcher.unhook,
		dispatcher.hook,
		dispatcher.oscDispatch,
		dispatcher.csiDispatch,
		dispatcher.escDispatch,
	)
	dispatcher.parser = parser
	return dispatcher
}

func (d *dispatcherImpl) Write(r io.Reader) error {
	defer func() {
		if err := recover(); err != nil {
			d.logger.Printf("[Write] panic: %v\n", err)
		}
	}()

	buff := make([]byte, 2048)

	for {
		n, err := r.Read(buff)

		if err != nil {
			if err == io.EOF {
				d.logger.Printf("[Write] EOF\n")
				return nil
			}

			d.logger.Printf("[Write] error: %v\n", err)
			return err
		}

		for _, b := range buff[:n] {
			d.parser.Advance(b)
		}
	}
}

func (d *dispatcherImpl) Flush() string {
	d.logger.Printf("[Flush] %s\n", string(d.cmd))
	str := string(d.cmd)
	d.cmd = []rune{}
	d.cur = 0
	return str
}

func (d *dispatcherImpl) Parse(r io.Reader) string {
	err := d.Write(r)
	if err != nil {
		return ""
	}
	return d.Flush()
}

func (d *dispatcherImpl) print(r rune) {
	d.logger.Printf("[Print] %c\n", r)
	if len(d.cmd) > d.cur {
		d.cmd[d.cur] = r
	} else {
		d.cmd = append(d.cmd, r)
	}
	d.cur++
}

func (d *dispatcherImpl) execute(b byte) {
	d.logger.Printf("[Execute] %02x\n", b)
	if b == '\b' {
		if d.cmd[d.cur-1] > 127 {
			d.cmd[d.cur-1] = ' '
		} else {
			d.cur--
		}
	}
}

func (d *dispatcherImpl) put(b byte) {
	fmt.Printf("[Put] %02x\n", b)
}

func (d *dispatcherImpl) unhook() {
	d.logger.Printf("[Unhook]\n")
}

func (d *dispatcherImpl) hook(params []int64, intermediates []byte, ignore bool, r rune) {
	d.logger.Printf("[Hook] params=%v, intermediates=%v, ignore=%v, r=%v\n", params, intermediates, ignore, r)
}

func (d *dispatcherImpl) oscDispatch(params [][]byte, bellTerminated bool) {
	d.logger.Printf("[OscDispatch] params=%v, bellTerminated=%v\n", params, bellTerminated)
}

func (d *dispatcherImpl) csiDispatch(params []int64, intermediates []byte, ignore bool, r rune) {
	d.logger.Printf("[CsiDispatch] params=%v, intermediates=%v, ignore=%v, r=%v\n", params, intermediates, ignore, r)
	switch r {
	case 64:
		d.cmd = append(append(d.cmd[:d.cur], 0), d.cmd[d.cur:]...)
	case 71:
		d.cur = int(params[0]) - 1
		d.cmd = d.cmd[:d.cur]
	case 75:
		if len(d.cmd) > d.cur {
			d.cmd = append(d.cmd[:d.cur], d.cmd[d.cur+1:]...)
		}
	case 80:
		if len(d.cmd) > d.cur {
			d.cmd = append(d.cmd[:d.cur], d.cmd[d.cur+1:]...)
		}
	default:
		d.logger.Printf("[CsiDispatch] unknown command: %v\n", r)
	}
}

func (d *dispatcherImpl) escDispatch(intermediates []byte, ignore bool, b byte) {
	d.logger.Printf("[EscDispatch] intermediates=%v, ignore=%v, byte=%02x\n", intermediates, ignore, b)
}
