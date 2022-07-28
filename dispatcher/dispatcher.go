package dispatcher

import (
	"fmt"
	"github.com/danielgatis/go-vte/vtparser"
	"io"
	"log"
)

var _ Dispatcher = &dispatcherImpl{}

var defaultDispatcher Dispatcher

// Dispatcher is the interface that wraps the methods of the dispatcher
type Dispatcher interface {
	Parse(r io.Reader) string
	Write(r io.Reader) error
	Flush() string
}

type dispatcherImpl struct {
	cmd []byte // command buffer
	cur int    // cursor position

	parser *vtparser.Parser

	logger *log.Logger
}

func Parse(r io.Reader) string {
	if defaultDispatcher == nil {
		defaultDispatcher = NewDispatcher()
	}
	return defaultDispatcher.Parse(r)
}

func NewDispatcher() Dispatcher {
	dispatcher := &dispatcherImpl{cmd: []byte{}, logger: log.New(io.Discard, "", log.LstdFlags)}
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
	dispatcher := &dispatcherImpl{cmd: []byte{}, logger: logger}
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
	d.logger.Printf("[Flush] %s\n", d.cmd)
	str := string(d.cmd)
	d.cmd = []byte{}
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
		d.cmd[d.cur] = byte(r)
	} else {
		d.cmd = append(d.cmd, byte(r))
	}
	d.cur++
}

func (d *dispatcherImpl) execute(b byte) {
	d.logger.Printf("[Execute] %02x\n", b)
	if b == '\b' {
		d.cur--
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
	case 71:
		d.cur = int(params[0]) - 1
		d.cmd = d.cmd[:d.cur]
	case 75:
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
