package trace

import (
	"fmt"
	"io"
)

// Tracer is the iterface that describes an object capable of tracing events throughout the code

type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Off crestes a tracer that will ignore calls to Trace
func Off() Tracer {
	return &nilTracer{}
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
