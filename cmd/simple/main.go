package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xyproto/midi"
)

func createChord(notes []string, startPause time.Duration) []midi.Note {
	var chord []midi.Note
	for _, note := range notes {
		chord = append(chord, midi.Note{
			Frequency:  midi.NoteNameToFrequency(note),
			Duration:   time.Second, // each note lasts for 1 second
			Velocity:   127,
			Channel:    1,
			StartPause: startPause,
		})
	}
	return chord
}

func main() {
	m := midi.NewMIDI(1, 480, 120) // Format 1 (multiple tracks), with 480 ticks per beat, and 120 BPM

	chords := [][]string{
		{"D4", "F4", "A4", "C5"},
		{"G4", "B4", "D5", "F5"},
		{"C4", "E4", "G4", "B4"},
		{"A4", "C5", "E5", "G5"},
	}

	for i, chordNotes := range chords {
		chord := createChord(chordNotes, time.Duration(i)*time.Second)
		for _, note := range chord {
			t := midi.NewTrack()
			m.AddTrack(t)
			m.AddNote(t, &note)
		}
	}

	f, err := os.Create("output.mid")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = m.Write(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MIDI file has been written to output.mid")
}
