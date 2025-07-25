package service

import (
	"fmt"
	"strings"
)

func NewDebugger(enabled bool) Debugger {
	var buf strings.Builder
	return debug{isEnabled: enabled, buf: &buf}
}

type debug struct {
	buf       *strings.Builder
	isEnabled bool
}

func (d debug) Write(a ...any) {
	if d.isEnabled {
		fmt.Fprint(d.buf, a...)
	}
}

func (d debug) Writeln(a ...any) {
	if d.isEnabled {
		fmt.Fprintln(d.buf, a...)
	}
}

func (d debug) Read() {
	if d.isEnabled {
		fmt.Print(d.buf.String())
	}
}
