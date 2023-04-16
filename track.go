package midi

import (
	"io"
	"time"
)

// ConvertToMIDITracks converts a slice of MidiNotes into a slice of MIDI tracks.
func ConvertToMIDITracks(notes []MidiNote) [][]byte {
	tracks := [][]byte{nil}
	lastNoteOff := -1
	lastNoteTime := 0
	for _, note := range notes {
		deltaTime := int(note.Duration/time.Millisecond) - lastNoteTime
		if deltaTime < 0 {
			deltaTime = 0
		}
		lastNoteTime = int(note.Duration / time.Millisecond)

		midiNote, pitchBend := FrequencyToMidi(note.Frequency)

		// If this note starts after the previous one, add a new track
		if deltaTime > 0 {
			tracks = append(tracks, nil)
			lastNoteOff = -1
		}

		// If this note starts after the last note ended, add a new event to the track
		if lastNoteOff != -1 {
			tracks[len(tracks)-1] = append(tracks[len(tracks)-1], deltaTimeToVLQ(deltaTime)...)
			WriteNoteOff(tracks[len(tracks)-1], note.Channel, midiNote, 0)
		}

		tracks[len(tracks)-1] = append(tracks[len(tracks)-1], deltaTimeToVLQ(0)...)
		WriteNoteOn(tracks[len(tracks)-1], note.Channel, midiNote, note.Velocity)
		tracks[len(tracks)-1] = append(tracks[len(tracks)-1], deltaTimeToVLQ(0)...)
		WritePitchBend(tracks[len(tracks)-1], note.Channel, pitchBend)

		lastNoteOff = len(tracks[len(tracks)-1]) - 3
	}
	return tracks
}

func WriteTrack(w io.Writer, notes []MidiNote) error {
	if len(notes) == 0 {
		return nil
	}

	deltaTime := 0
	var lastNoteOn MidiNote
	var err error

	for _, note := range notes {
		if note.Slur {
			if lastNoteOn.Frequency == note.Frequency && lastNoteOn.Channel == note.Channel {
				// continue playing the previous note
				continue
			} else {
				// end the previous note
				err = WriteNoteOff(w, lastNoteOn.Channel, lastNoteOn.MidiNote, deltaTime)
				if err != nil {
					return err
				}
				deltaTime = 0
			}
		} else {
			if lastNoteOn.Frequency != note.Frequency || lastNoteOn.Channel != note.Channel {
				// end the previous note if there was one
				if lastNoteOn.Frequency > 0 {
					err = WriteNoteOff(w, lastNoteOn.Channel, lastNoteOn.MidiNote, deltaTime)
					if err != nil {
						return err
					}
					deltaTime = 0
				}

				// start the new note
				midiNote, pitchBend := FrequencyToMidi(note.Frequency)
				err = WriteNoteOn(w, note.Channel, midiNote, note.Velocity)
				if err != nil {
					return err
				}
				err = WritePitchBend(w, note.Channel, pitchBend)
				if err != nil {
					return err
				}

				lastNoteOn = MidiNote{
					Frequency:  note.Frequency,
					MidiNote:   midiNote,
					Velocity:   note.Velocity,
					Channel:    note.Channel,
					Instrument: note.Instrument,
					Slur:       note.Slur,
				}

				deltaTime = 0
			}
		}

		deltaTime += int(DurationToTicks(note.Duration))
	}

	// end the last note
	if lastNoteOn.Frequency > 0 {
		err = WriteNoteOff(w, lastNoteOn.Channel, lastNoteOn.MidiNote, deltaTime)
		if err != nil {
			return err
		}
	}

	return nil
}
