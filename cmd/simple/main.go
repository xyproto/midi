package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xyproto/midi"
)

func main() {
	// Create a new MIDI object
	m := midi.NewMIDI(1, 480, 120) // Format 1 (multiple tracks), with 480 ticks per beat, and 120 BPM

	// Create a slice of MIDI notes
	notes := []midi.Note{
		{Frequency: midi.NoteToFrequency("C4"), Duration: time.Millisecond * 100, Velocity: 127, Channel: 1, Program: 1, StartPause: time.Millisecond * 0},
		{Frequency: midi.NoteToFrequency("D4"), Duration: time.Millisecond * 100, Velocity: 127, Channel: 1, Program: 1, StartPause: time.Millisecond * 0},
		{Frequency: midi.NoteToFrequency("E4"), Duration: time.Millisecond * 100, Velocity: 127, Channel: 1, Program: 1, StartPause: time.Millisecond * 0},
	}

	// Create a new track and add it to the MIDI object
	t := midi.NewTrack()
	m.AddTrack(t)

	// Add each note to the track
	for _, note := range notes {
		m.AddNote(t, &note)
	}

	// Write the MIDI object to a file
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
