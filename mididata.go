package midi

import (
	"encoding/binary"
	"io"
)

// This is a list of constants and utility functions for low-level MIDI operations.

// MIDI message status bytes
const (
	NoteOff               = 0x80
	NoteOn                = 0x90
	PolyphonicKeyPressure = 0xA0
	ControlChange         = 0xB0
	ProgramChange         = 0xC0
	ChannelPressure       = 0xD0
	PitchBend             = 0xE0
	SystemExclusive       = 0xF0
)

// writeVariableLengthQuantity writes a variable-length quantity (VLQ) to an io.Writer
func writeVariableLengthQuantity(w io.Writer, value uint32) error {
	bytes := make([]byte, 0, 4)
	bytes = append(bytes, byte(value&0x7F))
	value >>= 7
	for value > 0 {
		bytes = append(bytes, byte((value&0x7F)|0x80))
		value >>= 7
	}
	// Reverse the bytes, because we have generated them backwards
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	_, err := w.Write(bytes)
	return err
}

// readVariableLengthQuantity reads a variable-length quantity (VLQ) from an io.Reader
func readVariableLengthQuantity(r io.Reader) (uint32, error) {
	var value uint32
	var buf [1]byte

	for {
		if _, err := r.Read(buf[:]); err != nil {
			return 0, err
		}
		value = (value << 7) | uint32(buf[0]&0x7F)
		if buf[0]&0x80 == 0 {
			break
		}
	}

	return value, nil
}

// writeMIDIUint32 writes a big-endian uint32 to an io.Writer
func writeMIDIUint32(w io.Writer, value uint32) error {
	return binary.Write(w, binary.BigEndian, value)
}

// readMIDIUint32 reads a big-endian uint32 from an io.Reader
func readMIDIUint32(r io.Reader) (uint32, error) {
	var value uint32
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}

// writeMIDIUint16 writes a big-endian uint16 to an io.Writer
func writeMIDIUint16(w io.Writer, value uint16) error {
	return binary.Write(w, binary.BigEndian, value)
}

// readMIDIUint16 reads a big-endian uint16 from an io.Reader
func readMIDIUint16(r io.Reader) (uint16, error) {
	var value uint16
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}

// writeMIDIUint8 writes a uint8 to an io.Writer
func writeMIDIUint8(w io.Writer, value uint8) error {
	return binary.Write(w, binary.BigEndian, value)
}

// readMIDIUint8 reads a uint8 from an io.Reader
func readMIDIUint8(r io.Reader) (uint8, error) {
	var value uint8
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}

func uint16ToBytes(value uint16) []byte {
	return []byte{byte(value >> 8), byte(value & 0xFF)}
}

func uint32ToBytes(value uint32) []byte {
	return []byte{byte(value >> 24), byte((value >> 16) & 0xFF), byte((value >> 8) & 0xFF), byte(value & 0xFF)}
}

func writeUint32(w io.Writer, value uint32) error {
	return writeBytes(w, uint32ToBytes(value))
}

// writeBytes writes a byte slice to an io.Writer
func writeBytes(w io.Writer, data []byte) error {
	_, err := w.Write(data)
	return err
}
