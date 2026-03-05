package sqlitecloud

import (
	"fmt"
	"strings"
	"testing"
)

type testStringEnum string
type testIntEnum int

func TestProtocolBufferFromValueSupportsStringAlias(t *testing.T) {
	val := testStringEnum("active")
	buffers, err := protocolBufferFromValue(val)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(buffers) != 1 {
		t.Fatalf("expected 1 buffer, got %d", len(buffers))
	}
	got := string(buffers[0])
	want := fmt.Sprintf("%c%d %s\x00", CMD_ZEROSTRING, len("active")+1, "active")
	if got != want {
		t.Fatalf("unexpected encoded value: want %q got %q", want, got)
	}
}

func TestProtocolBufferFromValueSupportsIntAliasPointer(t *testing.T) {
	raw := testIntEnum(7)
	buffers, err := protocolBufferFromValue(&raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(buffers) != 1 {
		t.Fatalf("expected 1 buffer, got %d", len(buffers))
	}
	got := string(buffers[0])
	want := fmt.Sprintf("%c%d ", CMD_INT, 7)
	if got != want {
		t.Fatalf("unexpected encoded value: want %q got %q", want, got)
	}
}

func TestProtocolBufferFromValueSupportsFloat32(t *testing.T) {
	buffers, err := protocolBufferFromValue(float32(2.5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(buffers) != 1 {
		t.Fatalf("expected 1 buffer, got %d", len(buffers))
	}
	got := string(buffers[0])
	if !strings.HasPrefix(got, fmt.Sprintf("%c", CMD_FLOAT)) {
		t.Fatalf("expected float buffer prefix, got %q", got)
	}
}

func TestProtocolBufferFromValueUnsupportedTypeReturnsError(t *testing.T) {
	type unsupported struct {
		Name string
	}

	_, err := protocolBufferFromValue(unsupported{Name: "x"})
	if err == nil {
		t.Fatalf("expected error for unsupported type")
	}
}
