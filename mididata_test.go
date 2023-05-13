package midi

import (
	"bytes"
	"testing"
)

func TestFrequencyToMidi(t *testing.T) {
	note, bend := FrequencyToMidi(440.0)
	if note != 69 || bend != 0 {
		t.Errorf("FrequencyToMidi(440.0) = %d, %d, want 69, 0", note, bend)
	}

	note, bend = FrequencyToMidi(466.16)
	if note != 69 || bend <= 0 {
		t.Errorf("FrequencyToMidi(466.16) = %d, %d, want 69, positive bend", note, bend)
	}
}

func TestUint16ToBytes(t *testing.T) {
	bytes := uint16ToBytes(0x1234)
	if bytes[0] != 0x12 || bytes[1] != 0x34 {
		t.Errorf("uint16ToBytes(0x1234) = %X, want 12 34", bytes)
	}
}

func TestUint32ToBytes(t *testing.T) {
	bytes := uint32ToBytes(0x12345678)
	if bytes[0] != 0x12 || bytes[1] != 0x34 || bytes[2] != 0x56 || bytes[3] != 0x78 {
		t.Errorf("uint32ToBytes(0x12345678) = %X, want 12 34 56 78", bytes)
	}
}

func TestWriteUint32(t *testing.T) {
	buf := new(bytes.Buffer)
	err := writeUint32(buf, 0x12345678)
	if err != nil {
		t.Errorf("writeUint32 failed: %v", err)
	}
	bytes := buf.Bytes()
	if bytes[0] != 0x12 || bytes[1] != 0x34 || bytes[2] != 0x56 || bytes[3] != 0x78 {
		t.Errorf("writeUint32(0x12345678) = %X, want 12 34 56 78", bytes)
	}
}

func TestWriteMIDIUint32(t *testing.T) {
	buf := new(bytes.Buffer)
	err := writeMIDIUint32(buf, 0x12345678)
	if err != nil {
		t.Errorf("writeMIDIUint32 failed: %v", err)
	}
	bytes := buf.Bytes()
	if bytes[0] != 0x12 || bytes[1] != 0x34 || bytes[2] != 0x56 || bytes[3] != 0x78 {
		t.Errorf("writeMIDIUint32(0x12345678) = %X, want 12 34 56 78", bytes)
	}
}

func TestWriteAndReadMIDIUint32(t *testing.T) {
	buf := new(bytes.Buffer)
	err := writeMIDIUint32(buf, 0x12345678)
	if err != nil {
		t.Errorf("writeMIDIUint32 failed: %v", err)
	}
	value, err := readMIDIUint32(buf)
	if err != nil {
		t.Errorf("readMIDIUint32 failed: %v", err)
	}
	if value != 0x12345678 {
		t.Errorf("readMIDIUint32 = %X, want 12345678", value)
	}
}

func TestWriteAndReadMIDIUint16(t *testing.T) {
	buf := new(bytes.Buffer)
	err := writeMIDIUint16(buf, 0x1234)
	if err != nil {
		t.Errorf("writeMIDIUint16 failed: %v", err)
	}
	value, err := readMIDIUint16(buf)
	if err != nil {
		t.Errorf("readMIDIUint16 failed: %v", err)
	}
	if value != 0x1234 {
		t.Errorf("readMIDIUint16 = %X, want 1234", value)
	}
}

func TestWriteAndReadMIDIUint8(t *testing.T) {
	buf := new(bytes.Buffer)
	err := writeMIDIUint8(buf, 0x12)
	if err != nil {
		t.Errorf("writeMIDIUint8 failed: %v", err)
	}
	value, err := readMIDIUint8(buf)
	if err != nil {
		t.Errorf("readMIDIUint8 failed: %v", err)
	}
	if value != 0x12 {
		t.Errorf("readMIDIUint8 = %X, want 12", value)
	}
}

func TestWriteAndReadVariableLengthQuantity(t *testing.T) {
	buf := new(bytes.Buffer)
	err := writeVariableLengthQuantity(buf, 0x123456)
	if err != nil {
		t.Errorf("writeVariableLengthQuantity failed: %v", err)
	}
	value, err := readVariableLengthQuantity(buf)
	if err != nil {
		t.Errorf("readVariableLengthQuantity failed: %v", err)
	}
	if value != 0x123456 {
		t.Errorf("readVariableLengthQuantity = %X, want 123456", value)
	}
}
