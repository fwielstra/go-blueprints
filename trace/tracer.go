package trace

import (
	"fmt"
	"io"
)

//Tracer is an interface for easily tracking trace / logging information from the application.
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

//Trace will take the given values (similar to printf and co) and output them to the configured writer.
func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

//New creates a new Tracer instance.
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type nilTracer struct{}

//Trace the nilTracer's Trace function is a no-op.
func (t *nilTracer) Trace(a ...interface{}) {}

//Off Returns a no-op tracer.
func Off() Tracer {
	return &nilTracer{}
}
