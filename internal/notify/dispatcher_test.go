package notify

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

type stubChannel struct {
	name string
	err  error
	sent []alert.Event
}

func (s *stubChannel) Name() string { return s.name }
func (s *stubChannel) Send(e alert.Event) error {
	s.sent = append(s.sent, e)
	return s.err
}

func makeEvent() alert.Event {
	return alert.NewEvent(9999, "tcp", "unexpected listener")
}

func TestDispatch_AllChannelsCalled(t *testing.T) {
	a := &stubChannel{name: "a"}
	b := &stubChannel{name: "b"}
	d := NewDispatcher(nil, a, b)
	if err := d.Dispatch(makeEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a.sent) != 1 || len(b.sent) != 1 {
		t.Errorf("expected each channel to receive 1 event")
	}
}

func TestDispatch_PartialFailure(t *testing.T) {
	good := &stubChannel{name: "good"}
	bad := &stubChannel{name: "bad", err: errors.New("boom")}
	d := NewDispatcher(nil, good, bad)
	err := d.Dispatch(makeEvent())
	if err == nil {
		t.Fatal("expected error from failing channel")
	}
	if len(good.sent) != 1 {
		t.Errorf("good channel should still have received event")
	}
}

func TestDispatch_Add(t *testing.T) {
	d := NewDispatcher(nil)
	c := &stubChannel{name: "late"}
	d.Add(c)
	_ = d.Dispatch(makeEvent())
	if len(c.sent) != 1 {
		t.Errorf("expected dynamically added channel to receive event")
	}
}

func TestLogChannel_Send(t *testing.T) {
	var buf bytes.Buffer
	ch := NewLogChannel(&buf)
	ev := alert.Event{Port: 8080, Protocol: "tcp", Message: "test", Timestamp: time.Now()}
	if err := ch.Send(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "8080") || !strings.Contains(out, "tcp") {
		t.Errorf("log output missing expected fields: %q", out)
	}
}
