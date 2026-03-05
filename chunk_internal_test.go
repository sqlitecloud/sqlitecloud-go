package sqlitecloud

import (
	"fmt"
	"strings"
	"testing"
)

type testStringEnum string
type testIntEnum int

func TestProtocolBufferFromValue(t *testing.T) {
	type unsupported struct{}
	intVal := 42
	strVal := "hello"

	tests := []struct {
		name      string
		value     interface{}
		wantLen   int
		wantType  byte
		wantError bool
	}{
		{"nil", nil, 1, CMD_NULL, false},
		{"string", "hello", 1, CMD_ZEROSTRING, false},
		{"int", int(42), 1, CMD_INT, false},
		{"int8", int8(8), 1, CMD_INT, false},
		{"int16", int16(16), 1, CMD_INT, false},
		{"int32", int32(32), 1, CMD_INT, false},
		{"int64", int64(64), 1, CMD_INT, false},
		{"uint", uint(1), 1, CMD_INT, false},
		{"uint8", uint8(1), 1, CMD_INT, false},
		{"uint16", uint16(1), 1, CMD_INT, false},
		{"uint32", uint32(1), 1, CMD_INT, false},
		{"uint64", uint64(1), 1, CMD_INT, false},
		{"float32", float32(3.14), 1, CMD_FLOAT, false},
		{"float64", float64(2.71), 1, CMD_FLOAT, false},
		{"[]byte", []byte("blob"), 2, CMD_BLOB, false},
		{"bool true", true, 1, CMD_INT, false},
		{"bool false", false, 1, CMD_INT, false},
		{"*int", &intVal, 1, CMD_INT, false},
		{"*string", &strVal, 1, CMD_ZEROSTRING, false},
		{"*int nil", (*int)(nil), 1, CMD_NULL, false},
		{"*string nil", (*string)(nil), 1, CMD_NULL, false},
		{"unsupported", unsupported{}, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffers, err := protocolBufferFromValue(tt.value)
			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(buffers) != tt.wantLen {
				t.Fatalf("got %d buffers, want %d", len(buffers), tt.wantLen)
			}
			if tt.wantLen > 0 && buffers[0][0] != tt.wantType {
				t.Fatalf("got first buffer type %q, want %q", buffers[0][0], tt.wantType)
			}
		})
	}
}

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

func TestProtocolBufferFromValueMixedArrayNoSilentDrops(t *testing.T) {
	pInt := 99
	values := []interface{}{
		"hello",
		int(42),
		nil,
		&pInt,
		float64(3),
		uint(7),
		[]byte("x"),
		true,
	}

	buffers := [][]byte{}
	for i, v := range values {
		valueBuffers, err := protocolBufferFromValue(v)
		if err != nil {
			t.Fatalf("unexpected error at index %d (%T): %v", i, v, err)
		}
		if len(valueBuffers) == 0 {
			t.Fatalf("value at index %d produced zero buffers", i)
		}
		buffers = append(buffers, valueBuffers...)
	}

	if len(buffers) < len(values) {
		t.Fatalf("got %d total buffers, expected at least %d", len(buffers), len(values))
	}
}
