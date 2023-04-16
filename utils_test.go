package midi

import (
	"bytes"
	"testing"
)

func TestUint16ToBytes(t *testing.T) {
	tests := []struct {
		input    uint16
		expected []byte
	}{
		{0, []byte{0x00, 0x00}},
		{256, []byte{0x01, 0x00}},
		{65535, []byte{0xFF, 0xFF}},
	}

	for _, test := range tests {
		result := uint16ToBytes(test.input)
		if !bytes.Equal(result, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestUint32ToBytes(t *testing.T) {
	tests := []struct {
		input    uint32
		expected []byte
	}{
		{0, []byte{0x00, 0x00, 0x00, 0x00}},
		{65536, []byte{0x00, 0x01, 0x00, 0x00}},
		{4294967295, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, test := range tests {
		result := uint32ToBytes(test.input)
		if !bytes.Equal(result, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestWriteBytes(t *testing.T) {
	expected := []byte{0x01, 0x02, 0x03, 0x04}
	buf := &bytes.Buffer{}
	err := writeBytes(buf, expected)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	result := buf.Bytes()
	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestWriteUint32(t *testing.T) {
	tests := []struct {
		input    uint32
		expected []byte
	}{
		{0, []byte{0x00, 0x00, 0x00, 0x00}},
		{65536, []byte{0x00, 0x01, 0x00, 0x00}},
		{4294967295, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		err := writeUint32(buf, test.input)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		result := buf.Bytes()
		if !bytes.Equal(result, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}
