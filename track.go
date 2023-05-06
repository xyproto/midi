package midi

import (
	"io"
)

var lastNoteOn Note

func DeltaTimeLength(value int) int {
	if value < 0x80 {
		return 1
	}
	if value < 0x4000 {
		return 2
	}
	if value < 0x200000 {
		return 3
	}
	return 4
}

func WriteTrack(w io.Writer, notes []Note) error {
	// Calculate track length in bytes
	trackLength := 0
	for i, note := range notes {
		deltaTime := 0
		if i > 0 {
			deltaTime = int(note.Duration.Seconds() * 96.0)
		}
		trackLength += DeltaTimeLength(deltaTime)
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
			deltaTime = int(note.Duration.Seconds() * 96.0)
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
			noteOffDeltaTime := int(note.Duration.Seconds() * 96.0)
			WriteNoteOff(w, note.Channel, midiNote, noteOffDeltaTime)
		}

		lastNoteOn = note
	}

	return nil
}
