package midi

import (
	"bytes"
	"testing"
	"time"
)

func TestFrequencyToMidi(t *testing.T) {
	testCases := []struct {
		name      string
		frequency float64
		midiNote  int
		pitchBend int
	}{
		{"A4", 440.0, 69, 8192},
		{"C4", 261.63, 60, 8194},
		{"D5", 587.33, 74, 8192},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			midiNote, pitchBend := FrequencyToMidi(tc.frequency)

			if midiNote != tc.midiNote {
				t.Errorf("Expected MIDI note %d, got %d", tc.midiNote, midiNote)
			}

			if pitchBend != tc.pitchBend {
				t.Errorf("Expected pitch bend %d, got %d", tc.pitchBend, pitchBend)
			}
		})
	}
}

func TestConvertToMIDI(t *testing.T) {
	notes := []MidiNote{
		{Frequency: 440.0, Duration: time.Millisecond * 500, Velocity: 64, Channel: 1, Instrument: 1, Slur: false},
		{Frequency: 261.63, Duration: time.Millisecond * 500, Velocity: 64, Channel: 1, Instrument: 1, Slur: false},
		{Frequency: 587.33, Duration: time.Millisecond * 500, Velocity: 64, Channel: 1, Instrument: 1, Slur: false},
	}
	tracks := ConvertToMIDITracks(notes)

	midiData, err := ConvertToMIDI(tracks)
	if err != nil {
		t.Fatalf("Failed to convert to MIDI: %v", err)
	}

	if len(midiData) == 0 {
		t.Errorf("Expected non-empty MIDI data")
	}
}

func TestWriteTrack(t *testing.T) {
	notes := []MidiNote{
		{Frequency: 440.0, Duration: time.Millisecond * 500, Velocity: 64, Channel: 1, Instrument: 1, Slur: false},
		{Frequency: 261.63, Duration: time.Millisecond * 500, Velocity: 64, Channel: 1, Instrument: 1, Slur: false},
		{Frequency: 587.33, Duration: time.Millisecond * 500, Velocity: 64, Channel: 1, Instrument: 1, Slur: false},
	}

	buf := new(bytes.Buffer)
	err := WriteTrack(buf, notes)
	if err != nil {
		t.Fatalf("Failed to write track: %v", err)
	}

	if buf.Len() == 0 {
		t.Errorf("Expected non-empty track data")
	}
}
