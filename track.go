package midi

import (
	"bytes"
	"io"
)

var lastNoteOn MidiNote

func ConvertToMIDITracks(tracks [][]MidiNote) ([][]byte, error) {
	var midiTracks [][]byte
	for _, track := range tracks {
		buf := new(bytes.Buffer)

		// Write track header
		WriteHeader(buf, 1)

		// Write track data
		err := WriteTrack(buf, track)
		if err != nil {
			return nil, err
		}

		// Add track data to midiTracks
		midiTracks = append(midiTracks, buf.Bytes())
	}

	return midiTracks, nil
}

func WriteTrack(w io.Writer, notes []MidiNote) error {
	// Calculate track length in bytes
	trackLength := 0
	for i, note := range notes {
		deltaTime := 0
		if i > 0 {
			deltaTime = int(note.Duration.Seconds() * 96)
		}
		trackLength += deltaTimeLength(deltaTime)
		trackLength += 3 // Note on
		trackLength += 3 // Note off
	}

	// Write track header
	w.Write([]byte("MTrk"))
	w.Write(uint32ToBytes(uint32(trackLength)))

	// Write MIDI events
	for i, note := range notes {
		deltaTime := 0
		if i > 0 {
			deltaTime = int(note.Duration.Seconds() * 96)
		}

		// Write delta time
		err := writeDeltaTime(w, deltaTime)
		if err != nil {
			return err
		}

		// Write program change
		WriteProgramChange(w, note.Channel, note.Instrument)

		// Write note on
		midiNote, pitchBend := FrequencyToMidi(note.Frequency)
		WriteNoteOn(w, note.Channel, midiNote, note.Velocity)

		// Write pitch bend if needed
		if pitchBend != 8192 {
			WritePitchBend(w, note.Channel, pitchBend)
		}

		// Write note off
		if !note.Slur {
			WriteNoteOff(w, note.Channel, midiNote, deltaTime)
		}

		lastNoteOn = note
	}

	return nil
}
