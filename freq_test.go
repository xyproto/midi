package midi

import (
	"math"
	"testing"
)

func TestFrequencyToMidi(t *testing.T) {
	// Test for A4 (440Hz)
	midiNote, _ := FrequencyToMidi(440.0)
	if midiNote != 69 {
		t.Errorf("Expected MIDI note number for 440Hz is 69, but got %d", midiNote)
	}

	// Test for C4 (~261.63Hz)
	midiNote, _ = FrequencyToMidi(261.63)
	if midiNote != 60 {
		t.Errorf("Expected MIDI note number for 261.63Hz is 60, but got %d", midiNote)
	}

	// Test for D4 (~293.66Hz)
	midiNote, _ = FrequencyToMidi(293.66)
	if midiNote != 62 {
		t.Errorf("Expected MIDI note number for 293.66Hz is 62, but got %d", midiNote)
	}
}

func TestNoteToFrequency(t *testing.T) {
	// Test for A4
	frequency := NoteToFrequency("A4")
	if math.Abs(frequency-440.0) > 0.5 {
		t.Errorf("Expected frequency for A4 is 440Hz, but got %f", frequency)
	}

	// Test for C4
	frequency = NoteToFrequency("C4")
	if math.Abs(frequency-261.63) > 0.5 {
		t.Errorf("Expected frequency for C4 is 261.63Hz, but got %f", frequency)
	}

	// Test for D4
	frequency = NoteToFrequency("D4")
	if math.Abs(frequency-293.66) > 0.5 {
		t.Errorf("Expected frequency for D4 is 293.66Hz, but got %f", frequency)
	}
}
