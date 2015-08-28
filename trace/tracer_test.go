package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Return from new should not be nil")
	} else {
		tracer.Trace("Hello tracer package.")
		if buf.String() != "Hello tracer package.\n" {
			t.Errorf("Trace should not write: '%s' .", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("something..")
}
