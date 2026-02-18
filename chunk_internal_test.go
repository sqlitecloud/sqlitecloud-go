package sqlitecloud

import (
	"testing"
)

func TestProtocolBufferFromValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		wantLen  int // expected number of []byte buffers returned
		wantType byte
	}{
		// Basic types
		{"nil", nil, 1, CMD_NULL},
		{"string", "hello", 1, CMD_ZEROSTRING},
		{"int", int(42), 1, CMD_INT},
		{"int8", int8(8), 1, CMD_INT},
		{"int16", int16(16), 1, CMD_INT},
		{"int32", int32(32), 1, CMD_INT},
		{"int64", int64(64), 1, CMD_INT},
		{"float32", float32(3.14), 1, CMD_FLOAT},
		{"float64", float64(2.71), 1, CMD_FLOAT},
		{"[]byte", []byte("blob"), 2, CMD_BLOB}, // header + data

		// Unsigned integers
		{"uint", uint(1), 1, CMD_INT},
		{"uint8", uint8(1), 1, CMD_INT},
		{"uint16", uint16(1), 1, CMD_INT},
		{"uint32", uint32(1), 1, CMD_INT},
		{"uint64", uint64(1), 1, CMD_INT},

		// Pointer types (dereferenced)
		{"*int", intPtr(42), 1, CMD_INT},
		{"*string", strPtr("hello"), 1, CMD_ZEROSTRING},
		{"*int nil", (*int)(nil), 1, CMD_NULL},
		{"*string nil", (*string)(nil), 1, CMD_NULL},

		// Unsupported types still return empty buffers
		{"bool", true, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffers := protocolBufferFromValue(tt.value)
			if len(buffers) != tt.wantLen {
				t.Errorf("protocolBufferFromValue(%T(%v)): got %d buffers, want %d", tt.value, tt.value, len(buffers), tt.wantLen)
			}
			if tt.wantLen > 0 && len(buffers) > 0 {
				if buffers[0][0] != tt.wantType {
					t.Errorf("protocolBufferFromValue(%T(%v)): got type %c, want %c", tt.value, tt.value, buffers[0][0], tt.wantType)
				}
			}
		})
	}
}

func TestProtocolBufferFromValueMixedArray(t *testing.T) {
	// Simulates the loop in sendArray: builds buffers from a mixed values slice
	// and checks that the number of buffer groups matches the number of values.
	pInt := intPtr(99)
	values := []interface{}{
		"hello",     // string -> 1 buffer
		int(42),     // int -> 1 buffer
		nil,         // nil -> 1 buffer
		pInt,        // *int -> 1 buffer (dereferenced to int)
		float64(3),  // float64 -> 1 buffer
		uint(7),     // uint -> 1 buffer
		[]byte("x"), // []byte -> 2 buffers (header+data)
	}

	// Count how many values produce at least one buffer
	buffersPerValue := make([]int, len(values))
	totalBuffers := 0
	missingValues := 0

	for i, v := range values {
		bufs := protocolBufferFromValue(v)
		buffersPerValue[i] = len(bufs)
		totalBuffers += len(bufs)
		if len(bufs) == 0 {
			missingValues++
			t.Errorf("value[%d] (%T = %v) produced 0 buffers — will be silently dropped", i, v, v)
		}
	}

	if missingValues > 0 {
		t.Errorf("%d out of %d values produced no buffers and will be missing from the protocol message", missingValues, len(values))
	}

	// Reproduce the exact loop from sendArray
	buffers := [][]byte{}
	for _, v := range values {
		buffers = append(buffers, protocolBufferFromValue(v)...)
	}

	t.Logf("values count: %d, total buffers: %d, buffers per value: %v", len(values), len(buffers), buffersPerValue)

	// Every value must produce at least 1 buffer ([]byte produces 2)
	expectedMinBuffers := len(values)
	if len(buffers) < expectedMinBuffers {
		t.Errorf("buffers array has %d elements, expected at least %d (one per value). %d values were silently dropped.",
			len(buffers), expectedMinBuffers, missingValues)
	}
}

// helpers
func intPtr(v int) *int       { return &v }
func strPtr(v string) *string { return &v }
